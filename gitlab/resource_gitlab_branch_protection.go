package gitlab

import (
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

	schemaMergeExactlyOneOfList := []string{"merge_access_level", "allowed_to_merge"}
	schemaPushExactlyOneOfList := []string{"push_access_level", "allowed_to_push"}

	schemaAllowedToMergeAtLeastOneOfList := []string{"allowed_to_merge.0.user_id", "allowed_to_merge.0.group_id", "allowed_to_merge.0.access_level"}
	schemaAllowedToPushAtLeastOneOfList := []string{"allowed_to_push.0.user_id", "allowed_to_push.0.group_id", "allowed_to_push.0.access_level"}
	schemaAllowedToUnprotectAtLeastOneOfList := []string{"allowed_to_unprotect.0.user_id", "allowed_to_unprotect.0.group_id", "allowed_to_unprotect.0.access_level"}

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
				ExactlyOneOf: schemaMergeExactlyOneOfList,
				Optional:     true,
				ForceNew:     true,
			},
			"push_access_level": {
				Type:         schema.TypeString,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
				ExactlyOneOf: schemaPushExactlyOneOfList,
				Optional:     true,
				ForceNew:     true,
			},
			"unprotect_access_level": {
				Type:          schema.TypeString,
				ValidateFunc:  validateValueFunc(acceptedAccessLevels),
				Default:       accessLevel[gitlab.MaintainerPermissions],
				ConflictsWith: []string{"allowed_to_unprotect"},
				Optional:      true,
				ForceNew:      true,
			},
			"allowed_to_merge": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToMergeAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"group_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToMergeAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"access_level": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: schemaAllowedToMergeAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
					},
				},
				MaxItems:     1,
				ExactlyOneOf: schemaMergeExactlyOneOfList,
				Optional:     true,
				ForceNew:     true,
			},
			"allowed_to_push": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToPushAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"group_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToPushAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"access_level": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: schemaAllowedToPushAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
					},
				},
				MaxItems:     1,
				ExactlyOneOf: schemaPushExactlyOneOfList,
				Optional:     true,
				ForceNew:     true,
			},
			"allowed_to_unprotect": {
				Type: schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToUnprotectAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"group_id": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeInt,
							},
							AtLeastOneOf: schemaAllowedToUnprotectAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
						"access_level": {
							Type: schema.TypeList,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							AtLeastOneOf: schemaAllowedToUnprotectAtLeastOneOfList,
							Optional:     true,
							ForceNew:     true,
						},
					},
				},
				MaxItems:      1,
				ConflictsWith: []string{"unprotect_access_level"},
				Optional:      true,
				ForceNew:      true,
			},
			"code_owner_approval_required": {
				Type:     schema.TypeBool,
				Default:  false,
				Optional: true,
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
	unprotectAccessLevel := accessLevelID[d.Get("unprotect_access_level").(string)]
	allowedToMerge := d.Get("allowed_to_merge").([]interface{})
	allowedToPush := d.Get("allowed_to_push").([]interface{})
	allowedToUnprotect := d.Get("allowed_to_unprotect").([]interface{})
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	options := &gitlab.ProtectRepositoryBranchesOptions{
		Name:                      branch,
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
	}

	// Merge access
	if len(allowedToMerge) > 0 {
		allowedToMergePerms := constructProtectBranchPermissionOptions(allowedToMerge[0].(map[string]interface{}))
		options.AllowedToMerge = allowedToMergePerms
	} else {
		options.MergeAccessLevel = &mergeAccessLevel
	}

	// Push access
	if len(allowedToPush) > 0 {
		allowedToPushPerms := constructProtectBranchPermissionOptions(allowedToPush[0].(map[string]interface{}))
		options.AllowedToPush = allowedToPushPerms
	} else {
		options.PushAccessLevel = &pushAccessLevel
	}

	// Unprotect access
	if len(allowedToUnprotect) > 0 {
		allowedToUnprotectPerms := constructProtectBranchPermissionOptions(allowedToUnprotect[0].(map[string]interface{}))
		options.AllowedToUnprotect = allowedToUnprotectPerms
	} else {
		options.UnprotectAccessLevel = &unprotectAccessLevel
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

	return resourceGitlabBranchProtectionRead(d, meta)
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
	d.Set("merge_access_level", convertBranchAccessDescriptionsToStateBranchAccessDescription(pb.MergeAccessLevels))
	d.Set("push_access_level", convertBranchAccessDescriptionsToStateBranchAccessDescription(pb.PushAccessLevels))
	d.Set("unprotect_access_level", convertBranchAccessDescriptionsToStateBranchAccessDescription(pb.UnprotectAccessLevels))
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

func constructProtectBranchPermissionOptions(permsMap map[string]interface{}) []*gitlab.ProtectBranchPermissionOptions {
	var permissions []*gitlab.ProtectBranchPermissionOptions
	for key, value := range permsMap {
		for _, item := range value.([]interface{}) {
			var permission gitlab.ProtectBranchPermissionOptions
			switch key {
			case "user_id":
				userID := item.(int)
				permission.UserID = &userID
			case "group_id":
				groupID := item.(int)
				permission.GroupID = &groupID
			case "access_level":
				accessLevelID := accessLevelID[item.(string)]
				permission.AccessLevel = &accessLevelID
			}

			permissions = append(permissions, &permission)
		}
	}

	return permissions
}

type StateBranchAccessDescription struct {
	AccessLevel []string `json:"access_level"`
	GroupId     []int    `json:"group_id,omitempty"`
	UserId      []int    `json:"user_id,omitempty"`
}

func convertBranchAccessDescriptionsToStateBranchAccessDescription(descriptions []*gitlab.BranchAccessDescription) *StateBranchAccessDescription {
	stateDescription := StateBranchAccessDescription{}
	for _, description := range descriptions {
		stateDescription.AccessLevel = append(stateDescription.AccessLevel, accessLevel[description.AccessLevel])
		stateDescription.GroupId = append(stateDescription.GroupId, description.GroupID)
		stateDescription.UserId = append(stateDescription.UserId, description.UserID)
	}

	return &stateDescription
}
