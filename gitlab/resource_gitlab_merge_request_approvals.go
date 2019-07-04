package gitlab

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabMergeRequestApprovals() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))

	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Create: resourceGitlabMergeRequestApprovalsCreate,
		Read:   resourceGitlabMergeRequestApprovalsRead,
		Delete: resourceGitlabMergeRequestApprovalsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"approvals_before_merge": {
				Type:     schema.TypeInt,
				ForceNew: true,
				Required: true,
			},
			"reset_approvals_on_push": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"disable_overriding_approvers_per_merge_request": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
			"merge_requests_author_approval": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Optional: true,
			},
		},
	}
}

func resourceGitlabMergeRequestApprovalsCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	approvalsBeforeMerge := d.Get("approvals_before_merge").(int)
	resetApprovalsOnPush := d.Get("reset_approvals_on_push").(bool)
	disableOverridingApprovers := d.Get("disable_overriding_approvers_per_merge_request").(bool)
	mergeRequestsAuthorsApproval := d.Get("merge_requests_author_approval").(bool)

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ApprovalsBeforeMerge:                      &approvalsBeforeMerge,
		ResetApprovalsOnPush:                      &resetApprovalsOnPush,
		DisableOverridingApproversPerMergeRequest: &disableOverridingApprovers,
		MergeRequestsAuthorApproval:               &mergeRequestsAuthorsApproval,
	}

	log.Printf("[DEBUG] create gitlab approvals configuration for project %s", project)

	_, _, err := client.Projects.ChangeApprovalConfiguration(project, options)
	if err != nil {
		return err
	}

	d.SetId(strings.ReplaceAll(project, "/", ":"))
	return resourceGitlabMergeRequestApprovalsRead(d, meta)
}

func resourceGitlabMergeRequestApprovalsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := strings.ReplaceAll(d.Id(), ":", "/")

	log.Printf("[DEBUG] read gitlab approvals configuration for project %s", project)

	approvals, _, err := client.Projects.GetApprovalConfiguration(project)
	if err != nil {
		return err
	}

	d.Set("project", project)
	d.Set("approvalsBeforeMerge", approvals.ApprovalsBeforeMerge)
	d.Set("reset_approvals_on_push", approvals.ResetApprovalsOnPush)
	d.Set("disable_overriding_approvers_per_merge_request", approvals.DisableOverridingApproversPerMergeRequest)
	d.Set("merge_requests_author_approval", approvals.MergeRequestsAuthorApproval)

	return nil
}

func resourceGitlabMergeRequestApprovalsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := strings.ReplaceAll(d.Id(), ":", "/")

	log.Printf("[DEBUG] Delete gitlab approvals configuration for project %s", project)

	f := false
	zero := 0

	options := &gitlab.ChangeApprovalConfigurationOptions{
		ApprovalsBeforeMerge:                      &zero,
		ResetApprovalsOnPush:                      &f,
		DisableOverridingApproversPerMergeRequest: &f,
		MergeRequestsAuthorApproval:               &f,
	}
	_, _, err := client.Projects.ChangeApprovalConfiguration(project, options)
	return err
}
