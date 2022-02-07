package provider

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var (
	allowedToElem = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_level": {
				Description: "Level of access.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"access_level_description": {
				Description: "Readable description of level of access.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_id": {
				Description: "The ID of a GitLab user allowed to perform the relevant action. Mutually exclusive with `group_id`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
			"group_id": {
				Description: "The ID of a GitLab group allowed to perform the relevant action. Mutually exclusive with `user_id`.",
				Type:        schema.TypeInt,
				Optional:    true,
			},
		},
	}
)

var _ = registerResource("gitlab_branch_protection", func() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to protect a specific branch by an access level so that the user with less access level cannot Merge/Push to the branch.\n\n" +
			"-> The `allowed_to_push`, `allowed_to_merge` and `code_owner_approval_required` arguments require a GitLab Premium account or above.  Please refer to [Gitlab API documentation](https://docs.gitlab.com/ee/api/protected_branches.html) for further information.",

		CreateContext: resourceGitlabBranchProtectionCreate,
		ReadContext:   resourceGitlabBranchProtectionRead,
		UpdateContext: resourceGitlabBranchProtectionUpdate,
		DeleteContext: resourceGitlabBranchProtectionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The id of the project.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"branch": {
				Description: "Name of the branch.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"merge_access_level": {
				Description:      fmt.Sprintf("Access levels allowed to merge. Valid values are: %s.", renderValueListForDocs(validProtectedBranchTagAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validProtectedBranchTagAccessLevelNames, false)),
				Required:         true,
				ForceNew:         true,
			},
			"push_access_level": {
				Description:      fmt.Sprintf("Access levels allowed to push. Valid values are: %s.", renderValueListForDocs(validProtectedBranchTagAccessLevelNames)),
				Type:             schema.TypeString,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validProtectedBranchTagAccessLevelNames, false)),
				Required:         true,
				ForceNew:         true,
			},
			"allowed_to_push":  schemaAllowedTo(),
			"allowed_to_merge": schemaAllowedTo(),
			"code_owner_approval_required": {
				Description: "Can be set to true to require code owner approval before merging.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"branch_protection_id": {
				Description: "The ID of the branch protection (not the branch name).",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabBranchProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] create gitlab branch protection on branch %q for project %s", branch, project)

	if d.IsNewResource() {
		existing, resp, err := client.ProtectedBranches.GetProtectedBranch(project, branch, gitlab.WithContext(ctx))
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return diag.Errorf("error looking up protected branch %q on project %q: %v", branch, project, err)
		}
		if resp.StatusCode != http.StatusNotFound {
			return diag.Errorf("protected branch %q on project %q already exists: %+v", branch, project, *existing)
		}
	}

	mergeAccessLevel := accessLevelNameToValue[d.Get("merge_access_level").(string)]
	pushAccessLevel := accessLevelNameToValue[d.Get("push_access_level").(string)]
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	allowedToPush := expandBranchPermissionOptions(d.Get("allowed_to_push").(*schema.Set).List())
	allowedToMerge := expandBranchPermissionOptions(d.Get("allowed_to_merge").(*schema.Set).List())

	pb, _, err := client.ProtectedBranches.ProtectRepositoryBranches(project, &gitlab.ProtectRepositoryBranchesOptions{
		Name:                      &branch,
		PushAccessLevel:           &pushAccessLevel,
		MergeAccessLevel:          &mergeAccessLevel,
		AllowedToPush:             &allowedToPush,
		AllowedToMerge:            &allowedToMerge,
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
	}, gitlab.WithContext(ctx))
	if err != nil {
		return diag.Errorf("error protecting branch %q on project %q: %v", branch, project, err)
	}

	if !pb.CodeOwnerApprovalRequired && codeOwnerApprovalRequired {
		return diag.Errorf("feature unavailable: code owner approvals")
	}

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return resourceGitlabBranchProtectionRead(ctx, d, meta)
}

func resourceGitlabBranchProtectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, branch, err := projectAndBranchFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab branch protection for project %s, branch %s", project, branch)

	// Get protected branch by project ID/path and branch name
	pb, _, err := client.ProtectedBranches.GetProtectedBranch(project, branch, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to read gitlab branch protection for project %s, branch %s: %s", project, branch, err)
		d.SetId("")
		return nil
	}

	d.Set("project", project)
	d.Set("branch", pb.Name)

	pushAccessLevels := convertAllowedAccessLevelsToBranchAccessDescriptions(pb.PushAccessLevels)
	if len(pushAccessLevels) > 0 {
		if err := d.Set("push_access_level", pushAccessLevels[0].AccessLevel); err != nil {
			return diag.Errorf("error setting push_access_level: %v", err)
		}
	}

	mergeAccessLevels := convertAllowedAccessLevelsToBranchAccessDescriptions(pb.MergeAccessLevels)
	if len(mergeAccessLevels) > 0 {
		if err := d.Set("merge_access_level", mergeAccessLevels[0].AccessLevel); err != nil {
			return diag.Errorf("error setting merge_access_level: %v", err)
		}
	}

	// lintignore: R004 // TODO: Resolve this tfproviderlint issue
	if err := d.Set("allowed_to_push", convertAllowedToToBranchAccessDescriptions(pb.PushAccessLevels)); err != nil {
		return diag.Errorf("error setting allowed_to_push: %v", err)
	}
	// lintignore: R004 // TODO: Resolve this tfproviderlint issue
	if err := d.Set("allowed_to_merge", convertAllowedToToBranchAccessDescriptions(pb.MergeAccessLevels)); err != nil {
		return diag.Errorf("error setting allowed_to_merge: %v", err)
	}

	if err := d.Set("code_owner_approval_required", pb.CodeOwnerApprovalRequired); err != nil {
		return diag.Errorf("error setting code_owner_approval_required: %v", err)
	}

	d.Set("branch_protection_id", pb.ID)

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return nil
}

func resourceGitlabBranchProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if _, err := client.ProtectedBranches.RequireCodeOwnerApprovals(project, branch, options, gitlab.WithContext(ctx)); err != nil {
		// The user might be running a version of GitLab that does not support this feature.
		// We enhance the generic 404 error with a more informative message.
		if errResponse, ok := err.(*gitlab.ErrorResponse); ok && errResponse.Response.StatusCode == 404 {
			return diag.Errorf("feature unavailable: code owner approvals: %v", err)
		}

		return diag.FromErr(err)
	}

	return resourceGitlabBranchProtectionRead(ctx, d, meta)
}

func resourceGitlabBranchProtectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] Delete gitlab protected branch %s for project %s", branch, project)

	_, err := client.ProtectedBranches.UnprotectRepositoryBranches(project, branch, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func projectAndBranchFromID(id string) (string, string, error) {
	project, branch, err := parseTwoPartID(id)

	if err != nil {
		log.Printf("[WARN] cannot get branch protection id from input: %v", id)
	}
	return project, branch, err
}

func expandBranchPermissionOptions(allowedTo []interface{}) []*gitlab.BranchPermissionOptions {
	result := make([]*gitlab.BranchPermissionOptions, 0)
	for _, v := range allowedTo {
		opt := &gitlab.BranchPermissionOptions{}
		if userID, ok := v.(map[string]interface{})["user_id"]; ok && userID != 0 {
			opt.UserID = gitlab.Int(userID.(int))
		}
		if groupID, ok := v.(map[string]interface{})["group_id"]; ok && groupID != 0 {
			opt.GroupID = gitlab.Int(groupID.(int))
		}
		result = append(result, opt)
	}
	return result
}

func schemaAllowedTo() *schema.Schema {
	return &schema.Schema{
		Description: "Defines permissions for action.",
		Type:        schema.TypeSet,
		Optional:    true,
		ForceNew:    true,
		Elem:        allowedToElem,
	}
}

func convertAllowedAccessLevelsToBranchAccessDescriptions(descriptions []*gitlab.BranchAccessDescription) []stateBranchAccessDescription {
	result := make([]stateBranchAccessDescription, 0)

	for _, description := range descriptions {
		if description.UserID != 0 || description.GroupID != 0 {
			continue
		}
		result = append(result, stateBranchAccessDescription{
			AccessLevel:            accessLevelValueToName[description.AccessLevel],
			AccessLevelDescription: description.AccessLevelDescription,
		})
	}

	return result
}

func convertAllowedToToBranchAccessDescriptions(descriptions []*gitlab.BranchAccessDescription) []stateBranchAccessDescription {
	result := make([]stateBranchAccessDescription, 0)

	for _, description := range descriptions {
		if description.UserID == 0 && description.GroupID == 0 {
			continue
		}
		result = append(result, stateBranchAccessDescription{
			AccessLevel:            accessLevelValueToName[description.AccessLevel],
			AccessLevelDescription: description.AccessLevelDescription,
			UserID:                 description.UserID,
			GroupID:                description.GroupID,
		})
	}

	return result
}
