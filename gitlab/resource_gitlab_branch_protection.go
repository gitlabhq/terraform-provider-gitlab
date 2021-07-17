package gitlab

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var (
	allowedToElem = &schema.Resource{
		Schema: map[string]*schema.Schema{
			"access_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"access_level_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"group_id": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
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
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
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
			"allowed_to_push":  schemaAllowedTo(),
			"allowed_to_merge": schemaAllowedTo(),
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
	branch := d.Get("branch").(string)

	log.Printf("[DEBUG] create gitlab branch protection on branch %q for project %s", branch, project)

	if d.IsNewResource() {
		existing, resp, err := client.ProtectedBranches.GetProtectedBranch(project, branch)
		if err != nil && resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("error looking up protected branch %q on project %q: %v", branch, project, err)
		}
		if resp.StatusCode != http.StatusNotFound {
			return fmt.Errorf("protected branch %q on project %q already exists: %+v", branch, project, *existing)
		}
	}

	mergeAccessLevel := accessLevelID[d.Get("merge_access_level").(string)]
	pushAccessLevel := accessLevelID[d.Get("push_access_level").(string)]
	codeOwnerApprovalRequired := d.Get("code_owner_approval_required").(bool)

	allowedToPush := expandBranchPermissionOptions(d.Get("allowed_to_push").(*schema.Set).List())
	allowedToMerge := expandBranchPermissionOptions(d.Get("allowed_to_merge").(*schema.Set).List())

	pb, _, err := client.ProtectedBranches.ProtectRepositoryBranches(project, &gitlab.ProtectRepositoryBranchesOptions{
		Name:                      &branch,
		PushAccessLevel:           &pushAccessLevel,
		MergeAccessLevel:          &mergeAccessLevel,
		AllowedToPush:             allowedToPush,
		AllowedToMerge:            allowedToMerge,
		CodeOwnerApprovalRequired: &codeOwnerApprovalRequired,
	})
	if err != nil {
		return fmt.Errorf("error protecting branch %q on project %q: %v", branch, project, err)
	}

	if !pb.CodeOwnerApprovalRequired && codeOwnerApprovalRequired {
		return fmt.Errorf("feature unavailable: code owner approvals")
	}

	d.SetId(buildTwoPartID(&project, &pb.Name))

	return resourceGitlabBranchProtectionRead(d, meta)
}

func resourceGitlabBranchProtectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project, branch, err := projectAndBranchFromID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab branch protection for project %s, branch %s", project, branch)

	// Get protected branch by project ID/path and branch name
	pb, _, err := client.ProtectedBranches.GetProtectedBranch(project, branch)
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
			return fmt.Errorf("error setting push_access_level: %v", err)
		}
	}

	mergeAccessLevels := convertAllowedAccessLevelsToBranchAccessDescriptions(pb.MergeAccessLevels)
	if len(mergeAccessLevels) > 0 {
		if err := d.Set("merge_access_level", mergeAccessLevels[0].AccessLevel); err != nil {
			return fmt.Errorf("error setting merge_access_level: %v", err)
		}
	}

	// lintignore: R004 // TODO: Resolve this tfproviderlint issue
	if err := d.Set("allowed_to_push", convertAllowedToToBranchAccessDescriptions(pb.PushAccessLevels)); err != nil {
		return fmt.Errorf("error setting allowed_to_push: %v", err)
	}
	// lintignore: R004 // TODO: Resolve this tfproviderlint issue
	if err := d.Set("allowed_to_merge", convertAllowedToToBranchAccessDescriptions(pb.MergeAccessLevels)); err != nil {
		return fmt.Errorf("error setting allowed_to_merge: %v", err)
	}

	if err := d.Set("code_owner_approval_required", pb.CodeOwnerApprovalRequired); err != nil {
		return fmt.Errorf("error setting code_owner_approval_required: %v", err)
	}

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
		Type:     schema.TypeSet,
		Optional: true,
		ForceNew: true,
		Elem:     allowedToElem,
	}
}

type stateBranchAccessDescription struct {
	AccessLevel            string `mapstructure:"access_level"`
	AccessLevelDescription string `mapstructure:"access_level_description"`
	GroupID                int    `mapstructure:"group_id,omitempty"`
	UserID                 int    `mapstructure:"user_id,omitempty"`
}

func convertAllowedAccessLevelsToBranchAccessDescriptions(descriptions []*gitlab.BranchAccessDescription) []stateBranchAccessDescription {
	result := make([]stateBranchAccessDescription, 0)

	for _, description := range descriptions {
		if description.UserID != 0 || description.GroupID != 0 {
			continue
		}
		result = append(result, stateBranchAccessDescription{
			AccessLevel:            accessLevel[description.AccessLevel],
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
			AccessLevel:            accessLevel[description.AccessLevel],
			AccessLevelDescription: description.AccessLevelDescription,
			UserID:                 description.UserID,
			GroupID:                description.GroupID,
		})
	}

	return result
}
