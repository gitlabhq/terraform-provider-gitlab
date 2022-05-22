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

var milestoneStateToStateEvent = map[string]string{
	"active": "activate",
	"closed": "close",
}

var _ = registerResource("gitlab_project_milestone", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_milestone`" + ` resource allows to manage the lifecycle of a project milestone.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/milestones.html)`,

		CreateContext: resourceGitlabProjectMilestoneCreate,
		ReadContext:   resourceGitlabProjectMilestoneRead,
		UpdateContext: resourceGitlabProjectMilestoneUpdate,
		DeleteContext: resourceGitlabProjectMilestoneDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: gitlabProjectMilestoneGetSchema(),
	}
})

func resourceGitlabProjectMilestoneCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	title := d.Get("title").(string)

	options := &gitlab.CreateMilestoneOptions{
		Title: &title,
	}
	if description, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(description.(string))
	}
	if startDate, ok := d.GetOk("start_date"); ok {
		parsedStartDate, err := parseISO8601Date(startDate.(string))
		if err != nil {
			return diag.Errorf("Failed to parse start_date: %s. %v", startDate.(string), err)
		}
		options.StartDate = parsedStartDate
	}
	if dueDate, ok := d.GetOk("due_date"); ok {
		parsedDueDate, err := parseISO8601Date(dueDate.(string))
		if err != nil {
			return diag.Errorf("Failed to parse due_date: %s. %v", dueDate.(string), err)
		}
		options.DueDate = parsedDueDate
	}

	log.Printf("[DEBUG] create gitlab milestone in project %s with title %s", project, title)
	milestone, resp, err := client.Milestones.CreateMilestone(project, options, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[WARN] failed to create gitlab milestone in project %s with title %s (response %v)", project, title, resp)
		return diag.FromErr(err)
	}
	d.SetId(resourceGitLabProjectMilestoneBuildId(project, milestone.ID))

	updateOptions := gitlab.UpdateMilestoneOptions{}
	if stateEvent, ok := d.GetOk("state"); ok {
		updateOptions.StateEvent = gitlab.String(milestoneStateToStateEvent[stateEvent.(string)])
	}
	if updateOptions != (gitlab.UpdateMilestoneOptions{}) {
		_, _, err := client.Milestones.UpdateMilestone(project, milestone.ID, &updateOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.Errorf("Failed to update milestone ID %d in project %s right after creation: %v", milestone.ID, project, err)
		}
	}

	return resourceGitlabProjectMilestoneRead(ctx, d, meta)
}

func resourceGitlabProjectMilestoneRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, milestoneID, err := resourceGitLabProjectMilestoneParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab milestone in project %s with ID %d", project, milestoneID)
	milestone, resp, err := client.Milestones.GetMilestone(project, milestoneID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[WARN] recieved 404 for gitlab milestone ID %d in project %s, removing from state", milestoneID, project)
			d.SetId("")
			return nil
		}
		log.Printf("[WARN] failed to read gitlab milestone ID %d in project %s. Response %v", milestoneID, project, resp)
		return diag.FromErr(err)
	}

	stateMap := gitlabProjectMilestoneToStateMap(project, milestone)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabProjectMilestoneUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, milestoneID, err := resourceGitLabProjectMilestoneParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	options := &gitlab.UpdateMilestoneOptions{}
	if d.HasChange("title") {
		options.Title = gitlab.String(d.Get("title").(string))
	}
	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}
	if d.HasChange("start_date") {
		startDate := d.Get("start_date").(string)
		parsedStartDate, err := parseISO8601Date(startDate)
		if err != nil {
			return diag.Errorf("Failed to parse due_date: %s. %v", startDate, err)
		}
		options.StartDate = parsedStartDate
	}
	if d.HasChange("due_date") {
		dueDate := d.Get("due_date").(string)
		parsedDueDate, err := parseISO8601Date(dueDate)
		if err != nil {
			return diag.Errorf("Failed to parse due_date: %s. %v", dueDate, err)
		}
		options.DueDate = parsedDueDate
	}
	if d.HasChange("state") {
		options.StateEvent = gitlab.String(milestoneStateToStateEvent[d.Get("state").(string)])
	}

	log.Printf("[DEBUG] update gitlab milestone in project %s with ID %d", project, milestoneID)
	_, _, err = client.Milestones.UpdateMilestone(project, milestoneID, options, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[WARN] failed to update gitlab milestone in project %s with ID %d", project, milestoneID)
		return diag.FromErr(err)
	}

	return resourceGitlabProjectMilestoneRead(ctx, d, meta)
}

func resourceGitlabProjectMilestoneDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, milestoneID, err := resourceGitLabProjectMilestoneParseId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete gitlab milestone in project %s with ID %d", project, milestoneID)
	resp, err := client.Milestones.DeleteMilestone(project, milestoneID, gitlab.WithContext(ctx))
	if err != nil {
		log.Printf("[DEBUG] failed to delete gitlab milestone in project %s with ID %d. Response %v", project, milestoneID, resp)
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitLabProjectMilestoneParseId(id string) (string, int, error) {
	project, milestone, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	milestoneID, err := strconv.Atoi(milestone)
	if err != nil {
		return "", 0, err
	}

	return project, milestoneID, nil
}

func resourceGitLabProjectMilestoneBuildId(project string, milestoneID int) string {
	stringMilestoneID := fmt.Sprintf("%d", milestoneID)
	return buildTwoPartID(&project, &stringMilestoneID)
}
