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

var _ = registerResource("gitlab_group_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_group_hook` + "`" + ` resource allows to manage the lifecycle of a group hook.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#hooks)`,

		CreateContext: resourceGitlabGroupHookCreate,
		ReadContext:   resourceGitlabGroupHookRead,
		UpdateContext: resourceGitlabGroupHookUpdate,
		DeleteContext: resourceGitlabGroupHookDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: gitlabGroupHookSchema(),
	}
})

func resourceGitlabGroupHookCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	options := &gitlab.AddGroupHookOptions{
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
		SubGroupEvents:           gitlab.Bool(d.Get("subgroup_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if v, ok := d.GetOk("token"); ok {
		options.Token = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab group hook %q", *options.URL)

	hook, _, err := client.Groups.AddGroupHook(group, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGitlabGroupHookBuildID(group, hook.ID))
	d.Set("token", options.Token)

	return resourceGitlabGroupHookRead(ctx, d, meta)
}

func resourceGitlabGroupHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	group, hookID, err := resourceGitlabGroupHookParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] read gitlab group hook %s/%d", group, hookID)

	client := meta.(*gitlab.Client)
	hook, _, err := client.Groups.GetGroupHook(group, hookID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group hook not found %s/%d, removing from state", group, hookID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabGroupHookToStateMap(group, hook)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabGroupHookUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	group, hookID, err := resourceGitlabGroupHookParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*gitlab.Client)
	options := &gitlab.EditGroupHookOptions{
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
		SubGroupEvents:           gitlab.Bool(d.Get("subgroup_events").(bool)),
		EnableSSLVerification:    gitlab.Bool(d.Get("enable_ssl_verification").(bool)),
	}

	if d.HasChange("token") {
		options.Token = gitlab.String(d.Get("token").(string))
	}

	log.Printf("[DEBUG] update gitlab group hook %s", d.Id())

	_, _, err = client.Groups.EditGroupHook(group, hookID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabGroupHookRead(ctx, d, meta)
}

func resourceGitlabGroupHookDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	group, hookID, err := resourceGitlabGroupHookParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Delete gitlab group hook %s/%d", group, hookID)

	client := meta.(*gitlab.Client)
	_, err = client.Groups.DeleteGroupHook(group, hookID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabGroupHookBuildID(group string, agentID int) string {
	return fmt.Sprintf("%s:%d", group, agentID)
}

func resourceGitlabGroupHookParseID(id string) (string, int, error) {
	groupID, rawHookID, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	hookID, err := strconv.Atoi(rawHookID)
	if err != nil {
		return "", 0, err
	}

	return groupID, hookID, nil
}
