package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectHook() *schema.Resource {
	// lintignore: XR002 // TODO: Resolve this tfproviderlint issue
	return &schema.Resource{
		Create: resourceGitlabProjectHookCreate,
		Read:   resourceGitlabProjectHookRead,
		Update: resourceGitlabProjectHookUpdate,
		Delete: resourceGitlabProjectHookDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"push_events_branch_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"confidential_issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"merge_requests_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"tag_push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"note_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"confidential_note_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"job_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"pipeline_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"wiki_page_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"deployment_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"enable_ssl_verification": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGitlabProjectHookCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddProjectHookOptions{
		URL:                      gitlab.String(d.Get("url").(string)),
		PushEvents:               gitlab.Bool(d.Get("push_events").(bool)),
		PushEventsBranchFilter:   gitlab.String(d.Get("push_events_branch_filter").(string)),
		IssuesEvents:             gitlab.Bool(d.Get("issues_events").(bool)),
		ConfidentialIssuesEvents: gitlab.Bool(d.Get("confidential_issues_events").(bool)),
		MergeRequestsEvents:      gitlab.Bool(d.Get("merge_requests_events").(bool)),
		TagPushEvents:            gitlab.Bool(d.Get("tag_push_events").(bool)),
		NoteEvents:               gitlab.Bool(d.Get("note_events").(bool)),
		ConfidentialNoteEvents:   gitlab.Bool(d.Get("confidential_note_events").(bool)),
		JobEvents:                gitlab.Bool(d.Get("job_events").(bool)),
		PipelineEvents:           gitlab.Bool(d.Get("pipeline_events").(bool)),
		WikiPageEvents:           gitlab.Bool(d.Get("wiki_page_events").(bool)),
		DeploymentEvents:         gitlab.Bool(d.Get("deployment_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if v, ok := d.GetOk("token"); ok {
		options.Token = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab project hook %q", *options.URL)

	hook, _, err := client.Projects.AddProjectHook(project, options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", hook.ID))

	return resourceGitlabProjectHookRead(d, meta)
}

func resourceGitlabProjectHookRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] read gitlab project hook %s/%d", project, hookId)

	hook, _, err := client.Projects.GetProjectHook(project, hookId)
	if err != nil {
		return err
	}

	d.Set("url", hook.URL)
	d.Set("push_events", hook.PushEvents)
	d.Set("push_events_branch_filter", hook.PushEventsBranchFilter)
	d.Set("issues_events", hook.IssuesEvents)
	d.Set("confidential_issues_events", hook.ConfidentialIssuesEvents)
	d.Set("merge_requests_events", hook.MergeRequestsEvents)
	d.Set("tag_push_events", hook.TagPushEvents)
	d.Set("note_events", hook.NoteEvents)
	d.Set("confidential_note_events", hook.ConfidentialNoteEvents)
	d.Set("job_events", hook.JobEvents)
	d.Set("pipeline_events", hook.PipelineEvents)
	d.Set("wiki_page_events", hook.WikiPageEvents)
	d.Set("deployment_events", hook.DeploymentEvents)
	d.Set("enable_ssl_verification", hook.EnableSSLVerification)
	return nil
}

func resourceGitlabProjectHookUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	options := &gitlab.EditProjectHookOptions{
		URL:                      gitlab.String(d.Get("url").(string)),
		PushEvents:               gitlab.Bool(d.Get("push_events").(bool)),
		PushEventsBranchFilter:   gitlab.String(d.Get("push_events_branch_filter").(string)),
		IssuesEvents:             gitlab.Bool(d.Get("issues_events").(bool)),
		ConfidentialIssuesEvents: gitlab.Bool(d.Get("confidential_issues_events").(bool)),
		MergeRequestsEvents:      gitlab.Bool(d.Get("merge_requests_events").(bool)),
		TagPushEvents:            gitlab.Bool(d.Get("tag_push_events").(bool)),
		NoteEvents:               gitlab.Bool(d.Get("note_events").(bool)),
		ConfidentialNoteEvents:   gitlab.Bool(d.Get("confidential_note_events").(bool)),
		JobEvents:                gitlab.Bool(d.Get("job_events").(bool)),
		PipelineEvents:           gitlab.Bool(d.Get("pipeline_events").(bool)),
		WikiPageEvents:           gitlab.Bool(d.Get("wiki_page_events").(bool)),
		DeploymentEvents:         gitlab.Bool(d.Get("deployment_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if d.HasChange("token") {
		options.Token = gitlab.String(d.Get("token").(string))
	}

	log.Printf("[DEBUG] update gitlab project hook %s", d.Id())

	_, _, err = client.Projects.EditProjectHook(project, hookId, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectHookRead(d, meta)
}

func resourceGitlabProjectHookDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}
	log.Printf("[DEBUG] Delete gitlab project hook %s", d.Id())

	_, err = client.Projects.DeleteProjectHook(project, hookId)
	return err
}
