package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_hook", func() *schema.Resource {
	// lintignore: XR002 // TODO: Resolve this tfproviderlint issue
	return &schema.Resource{
		Description: "This resource allows you to create and manage hooks for your GitLab projects.\n" +
			"For further information on hooks, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ce/user/project/integrations/webhooks.html).",

		CreateContext: resourceGitlabProjectHookCreate,
		ReadContext:   resourceGitlabProjectHookRead,
		UpdateContext: resourceGitlabProjectHookUpdate,
		DeleteContext: resourceGitlabProjectHookDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the hook to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"url": {
				Description: "The url of the hook to invoke.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Description: "A token to present when invoking the hook.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
			},
			"push_events": {
				Description: "Invoke the hook for push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
			"push_events_branch_filter": {
				Description: "Invoke the hook for push events on matching branches only.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"issues_events": {
				Description: "Invoke the hook for issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"confidential_issues_events": {
				Description: "Invoke the hook for confidential issues events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"merge_requests_events": {
				Description: "Invoke the hook for merge requests.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"tag_push_events": {
				Description: "Invoke the hook for tag push events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"note_events": {
				Description: "Invoke the hook for notes events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"confidential_note_events": {
				Description: "Invoke the hook for confidential notes events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"job_events": {
				Description: "Invoke the hook for job events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"pipeline_events": {
				Description: "Invoke the hook for pipeline events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"wiki_page_events": {
				Description: "Invoke the hook for wiki page events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"deployment_events": {
				Description: "Invoke the hook for deployment events.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"enable_ssl_verification": {
				Description: "Enable ssl verification when invoking the hook.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
})

func resourceGitlabProjectHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	hook, _, err := client.Projects.AddProjectHook(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", hook.ID))

	return resourceGitlabProjectHookRead(ctx, d, meta)
}

func resourceGitlabProjectHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab project hook %s/%d", project, hookId)

	hook, _, err := client.Projects.GetProjectHook(project, hookId, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab project hook not found %s/%d", project, hookId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
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

func resourceGitlabProjectHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
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

	_, _, err = client.Projects.EditProjectHook(project, hookId, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabProjectHookRead(ctx, d, meta)
}

func resourceGitlabProjectHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookId, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab project hook %s", d.Id())

	_, err = client.Projects.DeleteProjectHook(project, hookId, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
