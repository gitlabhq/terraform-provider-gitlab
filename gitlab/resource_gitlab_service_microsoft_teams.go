package gitlab

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceMicrosoftTeams() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServiceMicrosoftTeamsCreate,
		Read:   resourceGitlabServiceMicrosoftTeamsRead,
		Update: resourceGitlabServiceMicrosoftTeamsUpdate,
		Delete: resourceGitlabServiceMicrosoftTeamsDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"active": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"webhook": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validateURLFunc,
			},
			"notify_only_broken_pipelines": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"branches_to_be_notified": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"push_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"confidential_issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"merge_requests_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"tag_push_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"note_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"confidential_note_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"pipeline_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"wiki_page_events": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourceGitlabServiceMicrosoftTeamsCreate(d *schema.ResourceData, meta interface{}) error {
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

	if _, err := client.Services.SetMicrosoftTeamsService(project, options); err != nil {
		return fmt.Errorf("couldn't create Gitlab Microsoft Teams service: %w", err)
	}

	return resourceGitlabServiceMicrosoftTeamsRead(d, meta)
}

func resourceGitlabServiceMicrosoftTeamsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	p, resp, err := client.Projects.GetProject(project, nil)
	if err != nil {
		if resp != nil && resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] Removing Gitlab Microsoft Teams service %s because project %s not found", d.Id(), p.Name)
			d.SetId("")
			return nil
		}
		return err
	}

	log.Printf("[DEBUG] Read Gitlab Microsoft Teams service for project %s", d.Id())

	teamsService, _, err := client.Services.GetMicrosoftTeamsService(project)
	if err != nil {
		return err
	}

	d.Set("webhook", teamsService.Properties.WebHook)
	d.Set("branches_to_be_notified", teamsService.Properties.BranchesToBeNotified)
	d.Set("notify_only_broken_pipelines", teamsService.Properties.NotifyOnlyBrokenPipelines)

	d.Set("project", project)
	d.Set("title", teamsService.Title)
	d.Set("created_at", teamsService.CreatedAt.String())
	d.Set("updated_at", teamsService.UpdatedAt.String())
	d.Set("active", teamsService.Active)
	d.Set("push_events", teamsService.PushEvents)
	d.Set("issues_events", teamsService.IssuesEvents)
	d.Set("merge_requests_events", teamsService.MergeRequestsEvents)
	d.Set("tag_push_events", teamsService.TagPushEvents)
	d.Set("note_events", teamsService.NoteEvents)
	d.Set("pipeline_events", teamsService.PipelineEvents)
	d.Set("confidential_issues_events", teamsService.ConfidentialIssuesEvents)
	d.Set("confidential_note_events", teamsService.ConfidentialNoteEvents)
	d.Set("wiki_page_events", teamsService.WikiPageEvents)

	return nil
}

func resourceGitlabServiceMicrosoftTeamsUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceGitlabServiceMicrosoftTeamsCreate(d, meta)
}

func resourceGitlabServiceMicrosoftTeamsDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Id()

	log.Printf("[DEBUG] Delete Gitlab Microsoft Teams service for project %s", d.Id())

	_, err := client.Services.DeleteMicrosoftTeamsService(project)
	return err
}
