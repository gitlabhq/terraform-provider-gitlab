package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var issueStateToStateEvent = map[string]string{
	"opened": "reopen",
	"closed": "close",
}

var _ = registerResource("gitlab_project_issue", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_issue`" + ` resource allows to manage the lifecycle of an issue within a project.

-> During a terraform destroy this resource will close the issue. Set the delete_on_destroy flag to true to delete the issue instead of closing it.

~> **Experimental** while the base functionality of this resource works, it may be subject to minor change.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/issues.html)
		`,

		CreateContext: resourceGitlabProjectIssueCreate,
		ReadContext:   resourceGitlabProjectIssueRead,
		UpdateContext: resourceGitlabProjectIssueUpdate,
		DeleteContext: resourceGitlabProjectIssueDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: constructSchema(
			gitlabProjectIssueGetSchema(),
			map[string]*schema.Schema{
				"delete_on_destroy": {
					Description: "Whether the issue is deleted instead of closed during destroy.",
					Type:        schema.TypeBool,
					Optional:    true,
					Default:     false,
				},
			},
		),
	}
})

func resourceGitlabProjectIssueCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)

	options := &gitlab.CreateIssueOptions{
		Title: gitlab.String(d.Get("title").(string)),
	}
	if iid, ok := d.GetOk("iid"); ok {
		options.IID = gitlab.Int(iid.(int))
	}
	if assigneeIDs, ok := d.GetOk("assignee_ids"); ok {
		options.AssigneeIDs = intSetToIntSlice(assigneeIDs.(*schema.Set))
	}
	if confidential, ok := d.GetOk("confidential"); ok {
		options.Confidential = gitlab.Bool(confidential.(bool))
	}
	if createdAt, ok := d.GetOk("created_at"); ok {
		parsedCreatedAt, err := time.Parse(time.RFC3339, createdAt.(string))
		if err != nil {
			return diag.Errorf("failed to parse created_at: %s. It must be in valid RFC3339 format.", err)
		}
		options.CreatedAt = gitlab.Time(parsedCreatedAt)
	}
	if description, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(description.(string))
	}
	if discussionToResolve, ok := d.GetOk("discussion_to_resolve"); ok {
		options.DiscussionToResolve = gitlab.String(discussionToResolve.(string))
	}
	if dueDate, ok := d.GetOk("due_date"); ok {
		parsedDueDate, err := parseISO8601Date(dueDate.(string))
		if err != nil {
			return diag.Errorf("failed to parse due_date: %s. %v", dueDate.(string), err)
		}
		options.DueDate = parsedDueDate
	}
	if issueType, ok := d.GetOk("issue_type"); ok {
		options.IssueType = gitlab.String(issueType.(string))
	}
	if labels, ok := d.GetOk("labels"); ok {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(labels.(*schema.Set)))
		options.Labels = &gitlabLabels
	}
	if mergeRequestToResolveDiscussionsOf, ok := d.GetOk("merge_request_to_resolve_discussions_of"); ok {
		options.MergeRequestToResolveDiscussionsOf = gitlab.Int(mergeRequestToResolveDiscussionsOf.(int))
	}
	if milestoneID, ok := d.GetOk("milestone_id"); ok {
		options.MilestoneID = gitlab.Int(milestoneID.(int))
	}
	if weight, ok := d.GetOk("weight"); ok {
		options.Weight = gitlab.Int(weight.(int))
	}

	issue, _, err := client.Issues.CreateIssue(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(resourceGitLabProjectIssueBuildId(project, issue.IID))

	updateOptions := gitlab.UpdateIssueOptions{}
	if discussionLocked, ok := d.GetOk("discussion_locked"); ok {
		updateOptions.DiscussionLocked = gitlab.Bool(discussionLocked.(bool))
	}
	if stateEvent, ok := d.GetOk("state"); ok {
		updateOptions.StateEvent = gitlab.String(issueStateToStateEvent[stateEvent.(string)])
	}
	if updateOptions != (gitlab.UpdateIssueOptions{}) {
		_, _, err := client.Issues.UpdateIssue(project, issue.IID, &updateOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("failed to update issue %d in project %s right after creation: %v", issue.IID, project, err)
		}
	}

	return resourceGitlabProjectIssueRead(ctx, d, meta)
}

func resourceGitlabProjectIssueRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueIID, err := resourceGitLabProjectIssueParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	issue, _, err := client.Issues.GetIssue(project, issueIID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[WARN] issue %d in project %s not found, removing from state", issueIID, project)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabProjectIssueToStateMap(project, issue)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabProjectIssueUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueIID, err := resourceGitLabProjectIssueParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.UpdateIssueOptions{}
	if d.HasChange("title") {
		options.Title = gitlab.String(d.Get("title").(string))
	}
	if d.HasChange("assignee_ids") {
		options.AssigneeIDs = intSetToIntSlice(d.Get("assignee_ids").(*schema.Set))
	}
	if d.HasChange("confidential") {
		options.Confidential = gitlab.Bool(d.Get("confidential").(bool))
	}
	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}
	if d.HasChange("due_date") {
		dueDate := d.Get("due_date").(string)
		if dueDate != "" {
			parsedDueDate, err := parseISO8601Date(dueDate)
			if err != nil {
				return diag.Errorf("failed to parse due_date: %s. %v", dueDate, err)
			}
			options.DueDate = parsedDueDate
		} else {
			// see https://github.com/xanzy/go-gitlab/issues/1384
			return diag.Errorf("remove a due date is currently not supported. See https://github.com/xanzy/go-gitlab/issues/1384")
		}
	}
	if d.HasChange("issue_type") {
		options.IssueType = gitlab.String(d.Get("issue_type").(string))
	}
	if d.HasChange("labels") {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(d.Get("labels").(*schema.Set)))
		options.Labels = &gitlabLabels
	}
	if d.HasChange("milestone_id") {
		options.MilestoneID = gitlab.Int(d.Get("milestone_id").(int))
	}
	if d.HasChange("weight") {
		options.Weight = gitlab.Int(d.Get("weight").(int))
	}
	if d.HasChange("state") {
		options.StateEvent = gitlab.String(issueStateToStateEvent[d.Get("state").(string)])
	}
	if d.HasChange("discussion_locked") {
		options.DiscussionLocked = gitlab.Bool(d.Get("discussion_locked").(bool))
	}

	_, _, err = client.Issues.UpdateIssue(project, issueIID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabProjectIssueRead(ctx, d, meta)
}

func resourceGitlabProjectIssueDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueIID, err := resourceGitLabProjectIssueParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	deleteOnDestroy := d.Get("delete_on_destroy").(bool)

	if deleteOnDestroy {
		log.Printf("[DEBUG] Deleting issue %d in project %s for destroy", issueIID, project)
		resp, err := client.Issues.DeleteIssue(project, issueIID, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("%s failed to delete issue %d in project %s: (%s) %v", d.Id(), issueIID, project, resp.Status, err)
		}
	} else {
		log.Printf("[DEBUG] Closing issue %d in project %s for destroy", issueIID, project)
		_, resp, err := client.Issues.UpdateIssue(project, issueIID, &gitlab.UpdateIssueOptions{StateEvent: gitlab.String("close")}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("%s failed to delete issue %d in project %s: (%s) %v", d.Id(), issueIID, project, resp.Status, err)
		}
	}

	return nil
}

func resourceGitLabProjectIssueParseId(id string) (string, int, error) {
	project, issue, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	issueIID, err := strconv.Atoi(issue)
	if err != nil {
		return "", 0, err
	}

	return project, issueIID, nil
}

func resourceGitLabProjectIssueBuildId(project string, issueIID int) string {
	stringIssueIID := fmt.Sprintf("%d", issueIID)
	return buildTwoPartID(&project, &stringIssueIID)
}
