package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectLevelMRApprovals() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectLevelMRApprovalsCreate,
		Read:   resourceGitlabProjectLevelMRApprovalsRead,
		Update: resourceGitlabProjectLevelMRApprovalsUpdate,
		Delete: resourceGitlabProjectLevelMRApprovalsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
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

func resourceGitlabProjectLevelMRApprovalsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	projectId := d.Get("project_id").(string)

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ResetApprovalsOnPush:                      gitlab.Bool(d.Get("reset_approvals_on_push").(bool)),
		DisableOverridingApproversPerMergeRequest: gitlab.Bool(d.Get("disable_overriding_approvers_per_merge_request").(bool)),
		MergeRequestsAuthorApproval:               gitlab.Bool(d.Get("merge_requests_author_approval").(bool)),
		MergeRequestsDisableCommittersApproval:    gitlab.Bool(d.Get("merge_requests_disable_committers_approval").(bool)),
	}

	log.Printf("[DEBUG] Creating new MR approval configuration for project %s: %#v", projectId, options)

	_, _, err := client.Projects.ChangeApprovalConfiguration(projectId, options)
	if err != nil {
		return fmt.Errorf("Error creating approval configuration: %s", err)
	}

	d.SetId(projectId)
	return resourceGitlabProjectLevelMRApprovalsRead(d, meta)
}

func resourceGitlabProjectLevelMRApprovalsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	projectId := d.Id()
	log.Printf("[DEBUG] Reading gitlab approval configuration for project %s", projectId)

	approvalConfig, _, err := client.Projects.GetApprovalConfiguration(projectId)
	if err != nil {
		return fmt.Errorf("Error reading approval configuration: %s", err)
	}

	d.Set("projectId", projectId)
	d.Set("reset_approvals_on_push", approvalConfig.ResetApprovalsOnPush)
	d.Set("disable_overriding_approvers_per_merge_request", approvalConfig.DisableOverridingApproversPerMergeRequest)
	d.Set("merge_requests_author_approval", approvalConfig.MergeRequestsAuthorApproval)
	d.Set("merge_requests_disable_committers_approval", approvalConfig.MergeRequestsDisableCommittersApproval)

	return nil
}

func resourceGitlabProjectLevelMRApprovalsUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.ChangeApprovalConfigurationOptions{}

	if d.HasChange("reset_approvals_on_push") {
		options.ResetApprovalsOnPush = gitlab.Bool(d.Get("reset_approvals_on_push").(bool))
	}
	if d.HasChange("disable_overriding_approvers_per_merge_request") {
		options.DisableOverridingApproversPerMergeRequest = gitlab.Bool(d.Get("disable_overriding_approvers_per_merge_request").(bool))
	}
	if d.HasChange("merge_requests_author_approval") {
		options.MergeRequestsAuthorApproval = gitlab.Bool(d.Get("merge_requests_author_approval").(bool))
	}
	if d.HasChange("merge_requests_disable_committers_approval") {
		options.MergeRequestsDisableCommittersApproval = gitlab.Bool(d.Get("merge_requests_disable_committers_approval").(bool))
	}

	_, _, err := client.Projects.ChangeApprovalConfiguration(d.Id(), options)
	if err != nil {
		return fmt.Errorf("Error updating approval configuration: %s", err)
	}

	return resourceGitlabProjectLevelMRApprovalsRead(d, meta)
}

func resourceGitlabProjectLevelMRApprovalsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	projectId := d.Id()

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ResetApprovalsOnPush:                      gitlab.Bool(true),
		DisableOverridingApproversPerMergeRequest: gitlab.Bool(false),
		MergeRequestsAuthorApproval:               gitlab.Bool(false),
		MergeRequestsDisableCommittersApproval:    gitlab.Bool(false),
	}

	log.Printf("[DEBUG] Resetting approval configuration for project %s: %#v", projectId, options)

	_, _, err := client.Projects.ChangeApprovalConfiguration(projectId, options)
	return err
}