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

var _ = registerResource("gitlab_project_issue_board", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_issue_board` + "`" + ` resource allows to manage the lifecycle of a Project Issue Board.

~> **NOTE:** If the board lists are changed all lists will be recreated.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/boards.html)`,

		CreateContext: resourceGitlabProjectIssueBoardCreate,
		ReadContext:   resourceGitlabProjectIssueBoardRead,
		UpdateContext: resourceGitlabProjectIssueBoardUpdate,
		DeleteContext: resourceGitlabProjectIssueBoardDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: gitlabProjectIssueBoardSchema(),
	}
})

func resourceGitlabProjectIssueBoardCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.CreateIssueBoardOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] create Project Issue Board %q in project %q", *options.Name, project)
	issueBoard, _, err := client.Boards.CreateIssueBoard(project, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGitlabProjectIssueBoardBuildID(project, issueBoard.ID))

	updateOptions := gitlab.UpdateIssueBoardOptions{}
	if v, ok := d.GetOk("milestone_id"); ok {
		updateOptions.MilestoneID = gitlab.Int(v.(int))
	}
	if v, ok := d.GetOk("assignee_id"); ok {
		updateOptions.AssigneeID = gitlab.Int(v.(int))
	}
	if v, ok := d.GetOk("labels"); ok {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(v.(*schema.Set)))
		updateOptions.Labels = &gitlabLabels
	}
	if v, ok := d.GetOk("weight"); ok {
		updateOptions.Weight = gitlab.Int(v.(int))
	}

	if (gitlab.UpdateIssueBoardOptions{}) != updateOptions {
		log.Printf("[DEBUG] update Project Issue Board %q in project %q after creation", *options.Name, project)
		_, _, err = client.Boards.UpdateIssueBoard(project, issueBoard.ID, &updateOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if v, ok := d.GetOk("lists"); ok {
		if err = resourceGitlabProjectIssueBoardCreateLists(ctx, client, project, issueBoard, v.([]interface{})); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabProjectIssueBoardRead(ctx, d, meta)
}

func resourceGitlabProjectIssueBoardRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueBoardID, err := resourceGitlabProjectIssueBoardParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read Project Issue Board in project %q with id %q", project, issueBoardID)
	issueBoard, _, err := client.Boards.GetIssueBoard(project, issueBoardID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] Project Issue Board in project %s with id %d not found, removing from state", project, issueBoardID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabProjectIssueBoardToStateMap(project, issueBoard)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabProjectIssueBoardUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueBoardID, err := resourceGitlabProjectIssueBoardParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.UpdateIssueBoardOptions{}
	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}
	if d.HasChange("milestone_id") {
		options.MilestoneID = gitlab.Int(d.Get("milestone_id").(int))
	}
	if d.HasChange("assignee_id") {
		options.AssigneeID = gitlab.Int(d.Get("assignee_id").(int))
	}
	if d.HasChange("labels") {
		gitlabLabels := gitlab.Labels(*stringSetToStringSlice(d.Get("labels").(*schema.Set)))
		options.Labels = &gitlabLabels
	}
	if d.HasChange("weight") {
		options.Weight = gitlab.Int(d.Get("weight").(int))
	}

	log.Printf("[DEBUG] update Project Issue Board %q in project %q", issueBoardID, project)
	updatedIssueBoard, _, err := client.Boards.UpdateIssueBoard(project, issueBoardID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("lists") {
		// NOTE: since we do not have a straightforward way to know which lists have been changed, we just re-create all lists
		log.Printf("[DEBUG] deleting lists for Project Issue Board %q in project %q", updatedIssueBoard.Name, project)
		for _, list := range updatedIssueBoard.Lists {
			log.Printf("[DEBUG] deleting list %d for Project Issue Board %q in project %q", list.ID, updatedIssueBoard.Name, project)
			_, err := client.Boards.DeleteIssueBoardList(project, issueBoardID, list.ID, gitlab.WithContext(ctx))
			if err != nil {
				return diag.Errorf("failed to delete list %q for Project Issue Board %q in project %q: %s", list.ID, updatedIssueBoard.Name, project, err)
			}
		}
		log.Printf("[DEBUG] deleted lists for Project Issue Board %q in project %q", updatedIssueBoard.Name, project)

		if err = resourceGitlabProjectIssueBoardCreateLists(ctx, client, project, updatedIssueBoard, d.Get("lists").([]interface{})); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabProjectIssueBoardRead(ctx, d, meta)
}

func resourceGitlabProjectIssueBoardDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, issueBoardID, err := resourceGitlabProjectIssueBoardParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete Project Issue Board in project %q with id %q", project, issueBoardID)
	if _, err := client.Boards.DeleteIssueBoard(project, issueBoardID, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabProjectIssueBoardBuildID(project string, issueBoardID int) string {
	return fmt.Sprintf("%s:%d", project, issueBoardID)
}

func resourceGitlabProjectIssueBoardParseID(id string) (string, int, error) {
	project, rawIssueBoardID, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	issueBoardID, err := strconv.Atoi(rawIssueBoardID)
	if err != nil {
		return "", 0, err
	}

	return project, issueBoardID, nil
}

func resourceGitlabProjectIssueBoardCreateLists(ctx context.Context, client *gitlab.Client, project string, issueBoard *gitlab.IssueBoard, lists []interface{}) error {
	log.Printf("[DEBUG] creating lists for Project Issue Board %q in project %q", issueBoard.Name, project)
	for i, listData := range lists {
		position := i + 1
		log.Printf("[DEBUG] creating list at position %d for Project Issue Board %q in project %q", position, issueBoard.Name, project)

		listOptions := gitlab.CreateIssueBoardListOptions{}
		if listData != nil {
			l := listData.(map[string]interface{})
			if v, ok := l["label_id"]; ok && v != 0 {
				listOptions.LabelID = gitlab.Int(v.(int))
			}
			if v, ok := l["assignee_id"]; ok && v != 0 {
				listOptions.AssigneeID = gitlab.Int(v.(int))
			}
			if v, ok := l["milestone_id"]; ok && v != 0 {
				listOptions.MilestoneID = gitlab.Int(v.(int))
			}
		}

		list, _, err := client.Boards.CreateIssueBoardList(project, issueBoard.ID, &listOptions, gitlab.WithContext(ctx))
		if err != nil {
			return fmt.Errorf("failed to create list at position %d for Project Issue Board %q in project %q: %s", position, issueBoard.Name, project, err)
		}

		log.Printf("[DEBUG] created list at position %d for Project Issue Board %q in project %q", list.Position, issueBoard.Name, project)
	}

	return nil
}
