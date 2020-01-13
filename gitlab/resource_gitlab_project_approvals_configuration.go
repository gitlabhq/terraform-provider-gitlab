package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectApprovalsConfiguration() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))

	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Create: resourceGitlabProjectApprovalsConfigurationCreate,
		Read:   resourceGitlabProjectApprovalsConfigurationRead,
		Update: resourceGitlabProjectApprovalsConfigurationUpdate,
		Delete: resourceGitlabProjectApprovalsConfigurationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"reset_approvals_on_push": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"disable_overriding_approvers_per_merge_request": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"merge_requests_author_approval": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"merge_requests_disable_committers_approval": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGitlabProjectApprovalsConfigurationCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	resetApprovalsOnPush := d.Get("reset_approvals_on_push").(bool)
	disableOverridingApproversPerMergeRequest := d.Get("disable_overriding_approvers_per_merge_request").(bool)
	mergeRequestsAuthorApproval := d.Get("merge_requests_author_approval").(bool)
	mergeRequestsDisableCommittersApproval := d.Get("merge_requests_disable_committers_approval").(bool)

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ResetApprovalsOnPush:                      gitlab.Bool(resetApprovalsOnPush),
		DisableOverridingApproversPerMergeRequest: gitlab.Bool(disableOverridingApproversPerMergeRequest),
		MergeRequestsAuthorApproval:               gitlab.Bool(mergeRequestsAuthorApproval),
		MergeRequestsDisableCommittersApproval:    gitlab.Bool(mergeRequestsDisableCommittersApproval),
	}

	log.Printf("[DEBUG] create gitlab approvals configuration for project %s, resetApprovalsOnPush: %t, disableOverridingApproversPerMergeRequest: %t, mergeRequestsAuthorApproval: %t, mergeRequestsDisableCommittersApproval: %t", project, , resetApprovalsOnPush, disableOverridingApproversPerMergeRequest, mergeRequestsAuthorApproval, mergeRequestsDisableCommittersApproval)

	_, _, err := client.Projects.ChangeApprovalConfiguration(project, options)
	if err != nil {
		return err
	}

	d.SetId(project)
	return resourceGitlabProjectApprovalsConfigurationRead(d, meta)
}

func resourceGitlabProjectApprovalsConfigurationRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	approvals, _, err := client.Projects.GetApprovalConfiguration(project)
	log.Printf("[DEBUG] read gitlab approvals configuration for project %s, resetApprovalsOnPush: %t, disableOverridingApproversPerMergeRequest: %t, mergeRequestsAuthorApproval: %t, mergeRequestsDisableCommittersApproval: %t", project, approvals.ResetApprovalsOnPush, approvals.DisableOverridingApproversPerMergeRequest, approvals.MergeRequestsAuthorApproval, approvals.MergeRequestsDisableCommittersApproval)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("reset_approvals_on_push", approvals.ResetApprovalsOnPush)
	d.Set("disable_overriding_approvers_per_merge_request", approvals.DisableOverridingApproversPerMergeRequest)
	d.Set("merge_requests_author_approval", approvals.MergeRequestsAuthorApproval)
	d.Set("merge_requests_disable_commiters_approval", approvals.MergeRequestsDisableCommittersApproval)

	return nil
}

func resourceGitlabProjectApprovalsConfigurationDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete (reset) gitlab approvals configuration for project %s", project)

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ResetApprovalsOnPush:                      gitlab.Bool(false),
		DisableOverridingApproversPerMergeRequest: gitlab.Bool(false),
		MergeRequestsAuthorApproval:               gitlab.Bool(false),
		MergeRequestsDisableCommittersApproval:    gitlab.Bool(false),
	}
	_, _, err := client.Projects.ChangeApprovalConfiguration(project, options)
	return err
}

func resourceGitlabProjectApprovalsConfigurationUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.ChangeApprovalConfigurationOptions{}
	if d.HasChange("reset_approvals_on_push") {
		options.ResetApprovalsOnPush = gitlab.Bool(d.Get("reset_approvals_on_push").(bool))
	}
	if d.HasChange("disable_overrading_approvers_per_merge_request") {
		options.DisableOverridingApproversPerMergeRequest = gitlab.Bool(d.Get("disable_overrading_approvers_per_merge_request").(bool))
	}
	if d.HasChange("merge_request_author_approval") {
		options.MergeRequestsAuthorApproval = gitlab.Bool(d.Get("merge_request_author_approval").(bool))
	}
	if d.HasChange("merge_request_disable_committers_approval") {
		options.MergeRequestsDisableCommittersApproval = gitlab.Bool(d.Get("merge_request_disable_committers_approval").(bool))
	}
	_, _, err := client.Projects.ChangeApprovalConfiguration(d.Id(), options)

	if err != nil {
		return err
	}
	return resourceGitlabProjectApprovalsConfigurationRead(d, meta)
}
