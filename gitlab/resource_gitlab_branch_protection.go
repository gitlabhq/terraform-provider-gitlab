package gitlab

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabBranchProtection() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))

	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Create: resourceGitlabBranchProtectionCreate,
		Read:   resourceGitlabBranchProtectionRead,
		Update: resourceGitlabBranchProtectionUpdate,
		Delete: resourceGitlabBranchProtectionDelete,
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"branch": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"merge_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
			"push_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				Required:     true,
				ForceNew:     true,
			},
			"users_allowed_to_merge": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"users_allowed_to_push": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
			"code_owner_approval_required": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGitlabBranchProtectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := gitlab.String(d.Get("branch").(string))
	mergeAccessLevel := accessLevelID[d.Get("merge_access_level").(string)]
	pushAccessLevel := accessLevelID[d.Get("push_access_level").(string)]
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	var allowedToMerge []*gitlab.ProtectBranchPermissionOptions
	usersAllowedToMerge := d.Get("users_allowed_to_merge")
	for _, userAllowedToMerge := range usersAllowedToMerge.(*schema.Set).List() {
		userID := userAllowedToMerge.(int)
		allowedToMerge = append(allowedToMerge, &gitlab.ProtectBranchPermissionOptions{UserID: &userID})
	}

	var allowedToPush []*gitlab.ProtectBranchPermissionOptions
	usersAllowedToPush := d.Get("users_allowed_to_push")
	for _, userAllowedToPush := range usersAllowedToPush.(*schema.Set).List() {
		userID := userAllowedToPush.(int)
		allowedToMerge = append(allowedToMerge, &gitlab.ProtectBranchPermissionOptions{UserID: &userID})
	}

	options := &gitlab.ProtectRepositoryBranchesOptions{
		Name:                      branch,
		MergeAccessLevel:          &mergeAccessLevel,
		PushAccessLevel:           &pushAccessLevel,
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
		AllowedToMerge:            allowedToMerge,
		AllowedToPush:             allowedToPush,
	}

	log.Printf("[DEBUG] create gitlab branch protection on %v for project %s", options.Name, project)

	bp, _, err := client.ProtectedBranches.ProtectRepositoryBranches(project, options)
	if err != nil {
		// Remove existing branch protection
		_, err = client.ProtectedBranches.UnprotectRepositoryBranches(project, *branch)
		if err != nil {
			return err
		}
		// Reprotect branch with updated values
		bp, _, err = client.ProtectedBranches.ProtectRepositoryBranches(project, options)
		if err != nil {
			return err
		}
	}

	d.SetId(buildTwoPartID(&project, &bp.Name))

	if err := resourceGitlabBranchProtectionRead(d, meta); err != nil {
		return err
	}

	// If the GitLab tier does not support the code owner approval feature, the resulting plan will be inconsistent.
	// We return an error because otherwise Terraform would report this inconsistency as a "bug in the provider" to the user.
	if codeOwnerApprovalRequired && !d.Get("code_owner_approval_required").(bool) {
		return errors.New("feature unavailable: code owner approvals")
	}

	return nil
}

func resourceGitlabBranchProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, branch, err := projectAndBranchFromID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab branch protection for project %s, branch %s", project, branch)

	pb, _, err := client.ProtectedBranches.GetProtectedBranch(project, branch)
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch protection for project %s, branch %s: %s", project, branch, err)
		d.SetId("")
		return nil
	}

	d.Set("project", project)
	d.Set("branch", pb.Name)
	var usersAllowedToMerge []int
	for _, mergeAccessLevel := range pb.MergeAccessLevels {
		if mergeAccessLevel.UserID == 0 && mergeAccessLevel.GroupID == 0 {
			d.Set("merge_access_level", accessLevel[mergeAccessLevel.AccessLevel])
			continue
		}
		if &mergeAccessLevel.UserID != nil {
			usersAllowedToMerge = append(usersAllowedToMerge, mergeAccessLevel.UserID)
		}
	}
	d.Set("users_allowed_to_merge", usersAllowedToMerge)
	var usersAllowedToPush []int
	for _, pushAccessLevel := range pb.PushAccessLevels {
		if pushAccessLevel.UserID == 0 && pushAccessLevel.GroupID == 0 {
			d.Set("push_access_level", accessLevel[pushAccessLevel.AccessLevel])
			continue
		}
		if &pushAccessLevel.UserID != nil {
			usersAllowedToPush = append(usersAllowedToPush, pushAccessLevel.UserID)
		}
	}
	d.Set("users_allowed_to_push", usersAllowedToPush)
	d.Set("code_owner_approval_required", pb.CodeOwnerApprovalRequired)

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return nil
}

func resourceGitlabBranchProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	// NOTE: At the time of writing, the only value that does not force re-creation is code_owner_approval_required,
	// so therefore that is the only update that needs to be handled.

	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	log.Printf("[DEBUG] update gitlab branch protection for project %s, branch %s", project, branch)

	options := &gitlab.RequireCodeOwnerApprovalsOptions{
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
	}

	if _, err := client.ProtectedBranches.RequireCodeOwnerApprovals(project, branch, options); err != nil {
		// The user might be running a version of GitLab that does not support this feature.
		// We enhance the generic 404 error with a more informative message.
		if errResponse, ok := err.(*gitlab.ErrorResponse); ok && errResponse.Response.StatusCode == 404 {
			return fmt.Errorf("feature unavailable: code owner approvals: %w", err)
		}

		return err
	}

	return resourceGitlabBranchProtectionRead(d, meta)
}

func resourceGitlabBranchProtectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] Delete gitlab protected branch %s for project %s", branch, project)

	_, err := client.ProtectedBranches.UnprotectRepositoryBranches(project, branch)
	return err
}

func projectAndBranchFromID(id string) (string, string, error) {
	project, branch, err := parseTwoPartID(id)

	if err != nil {
		log.Printf("[WARN] cannot get branch protection id from input: %v", id)
	}
	return project, branch, err
}
