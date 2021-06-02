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
			"allow_force_push": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
	allowForcePush := d.Get("allow_force_push").(bool)
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	options := &gitlab.ProtectRepositoryBranchesOptions{
		Name:                      branch,
		MergeAccessLevel:          &mergeAccessLevel,
		PushAccessLevel:           &pushAccessLevel,
		AllowForcePush:            &allowForcePush,
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
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
	d.Set("merge_access_level", accessLevel[pb.MergeAccessLevels[0].AccessLevel])
	d.Set("push_access_level", accessLevel[pb.PushAccessLevels[0].AccessLevel])
	d.Set("allow_force_push", pb.AllowForcePush)
	d.Set("code_owner_approval_required", pb.CodeOwnerApprovalRequired)

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return nil
}

func resourceGitlabBranchProtectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] update gitlab branch protection for project %s, branch %s", project, branch)

	if d.HasChange("allow_force_push") {
		allowForcePush := d.Get("allow_force_push").(bool)

		options := &gitlab.AllowForcePushOptions{
			AllowForcePush: &allowForcePush,
		}

		if _, err := client.ProtectedBranches.AllowForcePush(project, branch, options); err != nil {
			return err
		}
	}

	if d.HasChange("code_owner_approval_required") {
		codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

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
