package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_hook` + "`" + ` resource allows to manage the lifecycle of a project hook.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/projects.html#hooks)`,

		CreateContext: resourceGitlabProjectHookCreate,
		ReadContext:   resourceGitlabProjectHookRead,
		UpdateContext: resourceGitlabProjectHookUpdate,
		DeleteContext: resourceGitlabProjectHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabProjectHookStateImporter,
		},
		Schema: gitlabProjectHookSchema(),
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
		ReleasesEvents:           gitlab.Bool(d.Get("releases_events").(bool)),
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
	d.Set("token", options.Token)

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
			log.Printf("[DEBUG] gitlab project hook not found %s/%d, removing from state", project, hookId)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabProjectHookToStateMap(project, hook)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
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
		ReleasesEvents:           gitlab.Bool(d.Get("releases_events").(bool)),
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

func resourceGitlabProjectHookStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Project Hook import format; expected '{project_id}:{hook_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
