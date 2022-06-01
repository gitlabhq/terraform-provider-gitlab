package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_service_microsoft_teams", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_service_microsoft_teams`" + ` resource allows to manage the lifecycle of a project integration with Microsoft Teams.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/integrations.html#microsoft-teams)`,

		CreateContext: resourceGitlabServiceMicrosoftTeamsCreate,
		ReadContext:   resourceGitlabServiceMicrosoftTeamsRead,
		UpdateContext: resourceGitlabServiceMicrosoftTeamsUpdate,
		DeleteContext: resourceGitlabServiceMicrosoftTeamsDelete,
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
			"created_at": {
				Description: "Create time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"updated_at": {
				Description: "Update time.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"active": {
				Description: "Whether the integration is active.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"webhook": {
				Description:  "The Microsoft Teams webhook. For example, https://outlook.office.com/webhook/...",
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateURLFunc,
			},
			"notify_only_broken_pipelines": {
				Description: "Send notifications for broken pipelines",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"branches_to_be_notified": {
				Description: "Branches to send notifications for. Valid options are “all”, “default”, “protected”, and “default_and_protected”. The default value is “default”",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"push_events": {
				Description: "Enable notifications for push events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"issues_events": {
				Description: "Enable notifications for issue events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"confidential_issues_events": {
				Description: "Enable notifications for confidential issue events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"merge_requests_events": {
				Description: "Enable notifications for merge request events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"tag_push_events": {
				Description: "Enable notifications for tag push events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"note_events": {
				Description: "Enable notifications for note events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"confidential_note_events": {
				Description: "Enable notifications for confidential note events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"pipeline_events": {
				Description: "Enable notifications for pipeline events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"wiki_page_events": {
				Description: "Enable notifications for wiki page events",
				Type:        schema.TypeBool,
				Optional:    true,
			},
		},
	}
})

func resourceGitlabServiceMicrosoftTeamsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	d.SetId(project)

	options := &gitlab.SetMicrosoftTeamsServiceOptions{
		WebHook:                   gitlab.String(d.Get("webhook").(string)),
		NotifyOnlyBrokenPipelines: gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool)),
		BranchesToBeNotified:      gitlab.String(d.Get("branches_to_be_notified").(string)),
		PushEvents:                gitlab.Bool(d.Get("push_events").(bool)),
		IssuesEvents:              gitlab.Bool(d.Get("issues_events").(bool)),
		ConfidentialIssuesEvents:  gitlab.Bool(d.Get("confidential_issues_events").(bool)),
		MergeRequestsEvents:       gitlab.Bool(d.Get("merge_requests_events").(bool)),
		TagPushEvents:             gitlab.Bool(d.Get("tag_push_events").(bool)),
		NoteEvents:                gitlab.Bool(d.Get("note_events").(bool)),
		ConfidentialNoteEvents:    gitlab.Bool(d.Get("confidential_note_events").(bool)),
		PipelineEvents:            gitlab.Bool(d.Get("pipeline_events").(bool)),
		WikiPageEvents:            gitlab.Bool(d.Get("wiki_page_events").(bool)),
	}

	log.Printf("[DEBUG] Create Gitlab Microsoft Teams service")

	if _, err := client.Services.SetMicrosoftTeamsService(project, options, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("couldn't create Gitlab Microsoft Teams service: %v", err)
	}

	return resourceGitlabServiceMicrosoftTeamsRead(ctx, d, meta)
}

func resourceGitlabServiceMicrosoftTeamsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] Read Gitlab Microsoft Teams service for project %s", d.Id())

	teamsService, _, err := client.Services.GetMicrosoftTeamsService(project, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Unable to find Gitlab Microsoft Teams service in project %s, removing from state", project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("project", project)
	d.Set("created_at", teamsService.CreatedAt.String())
	d.Set("updated_at", teamsService.UpdatedAt.String())
	d.Set("active", teamsService.Active)
	d.Set("webhook", teamsService.Properties.WebHook)
	d.Set("notify_only_broken_pipelines", teamsService.Properties.NotifyOnlyBrokenPipelines)
	d.Set("branches_to_be_notified", teamsService.Properties.BranchesToBeNotified)
	d.Set("push_events", teamsService.PushEvents)
	d.Set("issues_events", teamsService.IssuesEvents)
	d.Set("confidential_issues_events", teamsService.ConfidentialIssuesEvents)
	d.Set("merge_requests_events", teamsService.MergeRequestsEvents)
	d.Set("tag_push_events", teamsService.TagPushEvents)
	d.Set("note_events", teamsService.NoteEvents)
	d.Set("confidential_note_events", teamsService.ConfidentialNoteEvents)
	d.Set("pipeline_events", teamsService.PipelineEvents)
	d.Set("wiki_page_events", teamsService.WikiPageEvents)

	return nil
}

func resourceGitlabServiceMicrosoftTeamsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceGitlabServiceMicrosoftTeamsCreate(ctx, d, meta)
}

func resourceGitlabServiceMicrosoftTeamsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] Delete Gitlab Microsoft Teams service for project %s", d.Id())

	_, err := client.Services.DeleteMicrosoftTeamsService(project, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return diag.FromErr(err)
}
