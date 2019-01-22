package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabServiceSlack() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabServiceSlackCreate,
		Read:   resourceGitlabServiceSlackRead,
		Update: resourceGitlabServiceSlackUpdate,
		Delete: resourceGitlabServiceSlackDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabServiceSlackImportState,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"webhook": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"notify_only_broken_pipelines": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"notify_only_default_branch": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"push_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"issue_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confidential_issues_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"confidential_issue_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"merge_requests_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"merge_request_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tag_push_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"tag_push_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"note_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"note_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"confidential_note_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			// TODO: Currently, GitLab ignores this option (not implemented yet?), so
			// there is no way to set it. Uncomment when this is fixed.
			// See: https://gitlab.com/gitlab-org/gitlab-ce/issues/49730
			//"confidential_note_channel": {
			//	Type:     schema.TypeString,
			//	Optional: true,
			//  Computed: true,
			//},
			"pipeline_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"pipeline_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"wiki_page_events": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"wiki_page_channel": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_events": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceGitlabServiceSlackSetToState(d *schema.ResourceData, service *gitlab.SlackService) {
	if !service.Active {
		log.Println("[DEBUG] slack service is not active")
		d.SetId("")
		return
	}

	d.SetId(fmt.Sprintf("%d", service.ID))
	d.Set("webhook", service.Properties.WebHook)
	d.Set("username", service.Properties.Username)
	d.Set("notify_only_broken_pipelines", service.Properties.NotifyOnlyBrokenPipelines.UnmarshalJSON)
	d.Set("notify_only_default_branch", service.Properties.NotifyOnlyDefaultBranch.UnmarshalJSON)
	d.Set("push_events", service.PushEvents)
	d.Set("push_channel", service.Properties.PushChannel)
	d.Set("issues_events", service.IssuesEvents)
	d.Set("issue_channel", service.Properties.IssueChannel)
	d.Set("confidential_issues_events", service.ConfidentialIssuesEvents)
	d.Set("confidential_issue_channel", service.Properties.ConfidentialIssueChannel)
	d.Set("merge_requests_events", service.MergeRequestsEvents)
	d.Set("merge_request_channel", service.Properties.MergeRequestChannel)
	d.Set("tag_push_events", service.TagPushEvents)
	d.Set("tag_push_channel", service.Properties.TagPushChannel)
	d.Set("note_events", service.NoteEvents)
	d.Set("note_channel", service.Properties.NoteChannel)
	d.Set("confidential_note_events", service.ConfidentialNoteEvents)
	// See comment to "confidential_note_channel" in resourceGitlabServiceSlack()
	//d.Set("confidential_note_channel", service.Properties.ConfidentialNoteChannel)
	d.Set("pipeline_events", service.PipelineEvents)
	d.Set("pipeline_channel", service.Properties.PipelineChannel)
	d.Set("wiki_page_events", service.WikiPageEvents)
	d.Set("wiki_page_channel", service.Properties.WikiPageChannel)
	d.Set("job_events", service.JobEvents)
}

func resourceGitlabServiceSlackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] create gitlab slack service for project %s", project)

	opts := &gitlab.SetSlackServiceOptions{
		WebHook: gitlab.String(d.Get("webhook").(string)),
	}

	if v := d.Get("username"); v != nil {
		opts.Username = gitlab.String(v.(string))
	}
	if v := d.Get("notify_only_broken_pipelines"); v != nil {
		opts.NotifyOnlyBrokenPipelines = gitlab.Bool(v.(bool))
	}
	if v := d.Get("notify_only_default_branch"); v != nil {
		opts.NotifyOnlyDefaultBranch = gitlab.Bool(v.(bool))
	}

	if v := d.Get("push_events"); v != nil {
		opts.PushEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("push_channel"); ok {
		opts.PushChannel = gitlab.String(v.(string))
	}

	if v := d.Get("issues_events"); v != nil {
		opts.IssuesEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("issue_channel"); ok {
		opts.IssueChannel = gitlab.String(v.(string))
	}

	if v := d.Get("confidential_issues_events"); v != nil {
		opts.ConfidentialIssuesEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("confidential_issue_channel"); ok {
		opts.ConfidentialIssueChannel = gitlab.String(v.(string))
	}

	if v := d.Get("merge_requests_events"); v != nil {
		opts.MergeRequestsEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("merge_request_channel"); ok {
		opts.MergeRequestChannel = gitlab.String(v.(string))
	}

	if v := d.Get("tag_push_events"); v != nil {
		opts.TagPushEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("tag_push_channel"); ok {
		opts.TagPushChannel = gitlab.String(v.(string))
	}

	if v := d.Get("note_events"); v != nil {
		opts.NoteEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("note_channel"); ok {
		opts.NoteChannel = gitlab.String(v.(string))
	}

	if v := d.Get("confidential_note_events"); v != nil {
		opts.ConfidentialNoteEvents = gitlab.Bool(v.(bool))
	}
	// See comment to "confidential_note_channel" in resourceGitlabServiceSlack()
	//if v, ok := d.GetOk("confidential_note_channel"); ok {
	//	opts.ConfidentialNoteChannel = gitlab.String(v.(string))
	//}

	if v := d.Get("pipeline_events"); v != nil {
		opts.PipelineEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("pipeline_channel"); ok {
		opts.PipelineChannel = gitlab.String(v.(string))
	}

	if v := d.Get("wiki_page_events"); v != nil {
		opts.WikiPageEvents = gitlab.Bool(v.(bool))
	}
	if v, ok := d.GetOk("wiki_page_channel"); ok {
		opts.WikiPageChannel = gitlab.String(v.(string))
	}

	_, err := client.Services.SetSlackService(project, opts)
	if err != nil {
		return err
	}

	return resourceGitlabServiceSlackRead(d, meta)
}

func resourceGitlabServiceSlackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] read gitlab slack service for project %s", project)

	service, _, err := client.Services.GetSlackService(project)
	if err != nil {
		return err
	}

	resourceGitlabServiceSlackSetToState(d, service)
	return nil
}

func resourceGitlabServiceSlackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] update gitlab slack service for project %s", project)

	opts := &gitlab.SetSlackServiceOptions{
		WebHook:  gitlab.String(d.Get("webhook").(string)),
		Username: gitlab.String(d.Get("username").(string)),
	}

	if d.HasChange("notify_only_broken_pipelines") {
		opts.NotifyOnlyBrokenPipelines = gitlab.Bool(d.Get("notify_only_broken_pipelines").(bool))
	}
	if d.HasChange("notify_only_default_branch") {
		opts.NotifyOnlyDefaultBranch = gitlab.Bool(d.Get("notify_only_default_branch").(bool))
	}

	if d.HasChange("push_events") {
		opts.PushEvents = gitlab.Bool(d.Get("push_events").(bool))
	}
	if d.HasChange("push_channel") {
		opts.PushChannel = gitlab.String(d.Get("push_channel").(string))
	}

	if d.HasChange("issues_events") {
		opts.IssuesEvents = gitlab.Bool(d.Get("issues_events").(bool))
	}
	if d.HasChange("issue_channel") {
		opts.IssueChannel = gitlab.String(d.Get("issue_channel").(string))
	}

	if d.HasChange("confidential_issues_events") {
		opts.ConfidentialIssuesEvents = gitlab.Bool(d.Get("confidential_issues_events").(bool))
	}
	if d.HasChange("confidential_issue_channel") {
		opts.ConfidentialIssueChannel = gitlab.String(d.Get("confidential_issue_channel").(string))
	}

	if d.HasChange("merge_requests_events") {
		opts.MergeRequestsEvents = gitlab.Bool(d.Get("merge_requests_events").(bool))
	}
	if d.HasChange("merge_request_channel") {
		opts.MergeRequestChannel = gitlab.String(d.Get("merge_request_channel").(string))
	}

	if d.HasChange("tag_push_events") {
		opts.TagPushEvents = gitlab.Bool(d.Get("tag_push_events").(bool))
	}
	if d.HasChange("tag_push_channel") {
		opts.TagPushChannel = gitlab.String(d.Get("tag_push_channel").(string))
	}

	if d.HasChange("note_events") {
		opts.NoteEvents = gitlab.Bool(d.Get("note_events").(bool))
	}
	if d.HasChange("note_channel") {
		opts.NoteChannel = gitlab.String(d.Get("note_channel").(string))
	}

	if d.HasChange("confidential_note_events") {
		opts.ConfidentialNoteEvents = gitlab.Bool(d.Get("confidential_note_events").(bool))
	}
	// See comment to "confidential_note_channel" in resourceGitlabServiceSlack()
	//if d.HasChange("confidential_note_channel") {
	//	opts.ConfidentialNoteChannel = gitlab.String(d.Get("confidential_note_channel").(string))
	//}

	if d.HasChange("pipeline_events") {
		opts.PipelineEvents = gitlab.Bool(d.Get("pipeline_events").(bool))
	}
	if d.HasChange("pipeline_channel") {
		opts.PipelineChannel = gitlab.String(d.Get("pipeline_channel").(string))
	}

	if d.HasChange("wiki_page_events") {
		opts.WikiPageEvents = gitlab.Bool(d.Get("wiki_page_events").(bool))
	}
	if d.HasChange("wiki_page_channel") {
		opts.WikiPageChannel = gitlab.String(d.Get("wiki_page_channel").(string))
	}

	_, err := client.Services.SetSlackService(project, opts)
	if err != nil {
		return err
	}

	return resourceGitlabServiceSlackRead(d, meta)
}

func resourceGitlabServiceSlackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	log.Printf("[DEBUG] delete gitlab slack service for project %s", project)

	_, err := client.Services.DeleteSlackService(project)
	return err
}
