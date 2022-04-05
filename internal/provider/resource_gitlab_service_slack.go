package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_service_slack", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_service_slack`" + ` resource allows to manage the lifecycle of a project integration with Slack.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/integrations.html#slack-notifications)`,

		CreateContext: resourceGitlabServiceSlackCreate,
		ReadContext:   resourceGitlabServiceSlackRead,
		UpdateContext: resourceGitlabServiceSlackUpdate,
		DeleteContext: resourceGitlabServiceSlackDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "ID of the project you want to activate integration on.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"webhook": {
				Description: "Webhook URL (ex.: https://hooks.slack.com/services/...)",
				Type:        schema.TypeString,
				Required:    true,
			},
			"username": {
				Description: "Username to use.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"notify_only_broken_pipelines": {
				Description: "Send notifications for broken pipelines.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"notify_only_default_branch": {
				Description: "This parameter has been replaced with `branches_to_be_notified`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Deprecated:  "use 'branches_to_be_notified' argument instead",
			},
			"branches_to_be_notified": {
				Description: "Branches to send notifications for. Valid options are \"all\", \"default\", \"protected\", and \"default_and_protected\".",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			// TODO: Currently, go-gitlab doesn't implement this option yet.
			//       see https://github.com/xanzy/go-gitlab/issues/1354
			// "commit_events": {
			// 	Description: "Enable notifications for commit events.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Computed:    true,
			// },
			"confidential_issue_channel": {
				Description: "The name of the channel to receive confidential issue events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"confidential_issues_events": {
				Description: "Enable notifications for confidential issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			// TODO: Currently, GitLab ignores this option (not implemented yet?), so
			// there is no way to set it. Uncomment when this is fixed.
			// See: https://gitlab.com/gitlab-org/gitlab-ce/issues/49730
			// "confidential_note_channel": {
			// 	Description: "The name of the channel to receive confidential note events notifications.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			"confidential_note_events": {
				Description: "Enable notifications for confidential note events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			// TODO: Currently, GitLab doesn't correctly implement the API, so this is
			//       impossible to implement here at the moment.
			//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
			// "deployment_channel": {
			// 	Description: "The name of the channel to receive deployment events notifications.",
			// 	Type:        schema.TypeString,
			// 	Optional:    true,
			// },
			// "deployment_events": {
			// 	Description: "Enable notifications for deployment events.",
			// 	Type:        schema.TypeBool,
			// 	Optional:    true,
			// 	Computed:    true,
			// },
			"issue_channel": {
				Description: "The name of the channel to receive issue events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"issues_events": {
				Description: "Enable notifications for issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			// TODO: Currently, go-gitlab doesn't implement this option yet.
			//       see https://github.com/xanzy/go-gitlab/issues/1354
			"job_events": {
				Description: "Enable notifications for job events. **ATTENTION**: This attribute is currently not being submitted to the GitLab API, due to https://github.com/xanzy/go-gitlab/issues/1354.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"merge_request_channel": {
				Description: "The name of the channel to receive merge request events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"merge_requests_events": {
				Description: "Enable notifications for merge requests events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"note_channel": {
				Description: "The name of the channel to receive note events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"note_events": {
				Description: "Enable notifications for note events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"pipeline_channel": {
				Description: "The name of the channel to receive pipeline events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"pipeline_events": {
				Description: "Enable notifications for pipeline events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"push_channel": {
				Description: "The name of the channel to receive push events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"push_events": {
				Description: "Enable notifications for push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"tag_push_channel": {
				Description: "The name of the channel to receive tag push events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"tag_push_events": {
				Description: "Enable notifications for tag push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
			"wiki_page_channel": {
				Description: "The name of the channel to receive wiki page events notifications.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"wiki_page_events": {
				Description: "Enable notifications for wiki page events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabServiceSlackCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)

	log.Printf("[DEBUG] create gitlab slack service for project %s", project)

	opts := &gitlab.SetSlackServiceOptions{
		WebHook: gitlab.String(d.Get("webhook").(string)),
	}

	opts.Username = gitlab.String(d.Get("username").(string))
	opts.NotifyOnlyBrokenPipelines = gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool))
	opts.NotifyOnlyDefaultBranch = gitlab.Bool(d.Get("notify_only_default_branch").(bool))
	opts.BranchesToBeNotified = gitlab.String(d.Get("branches_to_be_notified").(string))
	// TODO: Currently, go-gitlab doesn't implement this option yet.
	//       see https://github.com/xanzy/go-gitlab/issues/1354
	// opts.CommitEvents = gitlab.Bool(d.Get("commit_events").(bool))
	opts.ConfidentialIssueChannel = gitlab.String(d.Get("confidential_issue_channel").(string))
	opts.ConfidentialIssuesEvents = gitlab.Bool(d.Get("confidential_issues_events").(bool))
	// TODO: Currently, GitLab ignores this option (not implemented yet?), so
	// there is no way to set it. Uncomment when this is fixed.
	// See: https://gitlab.com/gitlab-org/gitlab-ce/issues/49730
	// opts.ConfidentialNoteChannel = gitlab.String(d.Get("confidential_note_channel").(string))
	opts.ConfidentialNoteEvents = gitlab.Bool(d.Get("confidential_note_events").(bool))
	// TODO: Currently, GitLab doesn't correctly implement the API, so this is
	//       impossible to implement here at the moment.
	//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
	// opts.DeploymentChannel = gitlab.String(d.Get("deployment_channel").(string))
	// opts.DeploymentEvents = gitlab.Bool(d.Get("deployment_events").(bool))
	opts.IssueChannel = gitlab.String(d.Get("issue_channel").(string))
	opts.IssuesEvents = gitlab.Bool(d.Get("issues_events").(bool))
	// TODO: Currently, go-gitlab doesn't implement this option yet.
	//       see https://github.com/xanzy/go-gitlab/issues/1354
	// opts.JobEvents = gitlab.Bool(d.Get("job_events").(bool))
	opts.MergeRequestChannel = gitlab.String(d.Get("merge_request_channel").(string))
	opts.MergeRequestsEvents = gitlab.Bool(d.Get("merge_requests_events").(bool))
	opts.NoteChannel = gitlab.String(d.Get("note_channel").(string))
	opts.NoteEvents = gitlab.Bool(d.Get("note_events").(bool))
	opts.PipelineChannel = gitlab.String(d.Get("pipeline_channel").(string))
	opts.PipelineEvents = gitlab.Bool(d.Get("pipeline_events").(bool))
	opts.PushChannel = gitlab.String(d.Get("push_channel").(string))
	opts.PushEvents = gitlab.Bool(d.Get("push_events").(bool))
	opts.TagPushChannel = gitlab.String(d.Get("tag_push_channel").(string))
	opts.TagPushEvents = gitlab.Bool(d.Get("tag_push_events").(bool))
	opts.WikiPageChannel = gitlab.String(d.Get("wiki_page_channel").(string))
	opts.WikiPageEvents = gitlab.Bool(d.Get("wiki_page_events").(bool))

	_, err := client.Services.SetSlackService(project, opts, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabServiceSlackRead(ctx, d, meta)
}

func resourceGitlabServiceSlackRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	id := d.Id()
	project := d.Get("project").(string)
	if id != project && project != "" {
		d.SetId(project)
		log.Printf("[WARN] changed gitlab slack service ID from %s to its project ID %s", id, project)
	} else {
		project = id
	}

	log.Printf("[DEBUG] read gitlab slack service for project %s", project)

	service, _, err := client.Services.GetSlackService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab slack service not found %s", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("webhook", service.Properties.WebHook)
	d.Set("username", service.Properties.Username)
	d.Set("notify_only_broken_pipelines", bool(service.Properties.NotifyOnlyBrokenPipelines))
	d.Set("notify_only_default_branch", bool(service.Properties.NotifyOnlyDefaultBranch))
	d.Set("branches_to_be_notified", service.Properties.BranchesToBeNotified)
	d.Set("confidential_issue_channel", service.Properties.ConfidentialIssueChannel)
	d.Set("confidential_issues_events", service.ConfidentialIssuesEvents)
	// TODO: Currently, GitLab ignores this option (not implemented yet?), so
	// there is no way to set it. Uncomment when this is fixed.
	// See: https://gitlab.com/gitlab-org/gitlab-ce/issues/49730
	// d.Set("confidential_note_channel", service.Properties.ConfidentialNoteChannel)
	d.Set("confidential_note_events", service.ConfidentialNoteEvents)
	// TODO: Currently, GitLab doesn't correctly implement the API, so this is
	//       impossible to implement here at the moment.
	//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
	// d.Set("deployment_channel", service.Properties.DeploymentChannel)
	// d.Set("deployment_events", service.DeploymentEvents)
	d.Set("issue_channel", service.Properties.IssueChannel)
	d.Set("issues_events", service.IssuesEvents)
	// TODO: Currently, go-gitlab doesn't implement this option yet.
	//       see https://github.com/xanzy/go-gitlab/issues/1354
	d.Set("job_events", service.JobEvents)
	d.Set("merge_request_channel", service.Properties.MergeRequestChannel)
	d.Set("merge_requests_events", service.MergeRequestsEvents)
	d.Set("note_channel", service.Properties.NoteChannel)
	d.Set("note_events", service.NoteEvents)
	d.Set("pipeline_channel", service.Properties.PipelineChannel)
	d.Set("pipeline_events", service.PipelineEvents)
	d.Set("push_channel", service.Properties.PushChannel)
	d.Set("push_events", service.PushEvents)
	d.Set("tag_push_channel", service.Properties.TagPushChannel)
	d.Set("tag_push_events", service.TagPushEvents)
	d.Set("wiki_page_channel", service.Properties.WikiPageChannel)
	d.Set("wiki_page_events", service.WikiPageEvents)

	return nil
}

func resourceGitlabServiceSlackUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGitlabServiceSlackCreate(ctx, d, meta)
}

func resourceGitlabServiceSlackDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] delete gitlab slack service for project %s", project)

	_, err := client.Services.DeleteSlackService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
