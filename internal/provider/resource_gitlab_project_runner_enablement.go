package provider

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_project_runner_enablement", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_runner_enablement`" + ` resource allows to enable a runner in a project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/runners.html#enable-a-runner-in-project)`,
		CreateContext: resourceGitlabProjectRunnerEnablementCreate,
		ReadContext:   resourceGitlabProjectRunnerEnablementRead,
		DeleteContext: resourceGitlabProjectRunnerEnablementDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"runner_id": {
				Description: "The ID of a runner to enable for the project.",
				Type:        schema.TypeInt,
				ForceNew:    true,
				Required:    true,
			},
		},
	}
})

func resourceGitlabProjectRunnerEnablementCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID := d.Get("project").(string)
	runnerID := d.Get("runner_id").(int)
	options := &gitlab.EnableProjectRunnerOptions{
		RunnerID: runnerID,
	}

	log.Printf("[DEBUG] create gitlab project runner %v/%v", projectID, runnerID)

	_, _, err := client.Runners.EnableProjectRunner(projectID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	runnerIDString := strconv.Itoa(runnerID)
	d.SetId(buildTwoPartID(&projectID, &runnerIDString))

	return resourceGitlabProjectRunnerEnablementRead(ctx, d, meta)
}

func resourceGitlabProjectRunnerEnablementRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, runnerID, err := projectAndRunnerFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab project runner %s/%v", project, runnerID)

	// Get the project id from `project`, which can be either the numeric ID or a name
	projectDetails, _, err := client.Projects.GetProject(project, &gitlab.GetProjectOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	runnerdetails, _, err := client.Runners.GetRunnerDetails(runnerID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	// Check if the project exists in the runner details
	found := false
	for _, p := range runnerdetails.Projects {
		if p.ID == projectDetails.ID {
			found = true
			break
		}
	}

	if !found {
		log.Printf("[WARN] removing project runner: %v from state because it no longer exists in gitlab", runnerID)
		d.SetId("")
		return nil
	}

	d.Set("project", project)
	d.Set("runner_id", runnerID)

	return nil
}

func projectAndRunnerFromID(id string) (string, int, error) {
	var runnerID int
	projectID, runnerIDString, err := parseTwoPartID(id)
	if err != nil {
		log.Printf("[WARN] could not get project and runner ids from resource id %v", id)
		return projectID, runnerID, err
	}

	runnerID, err = strconv.Atoi(runnerIDString)
	if err != nil {
		log.Printf("[WARN] could not convert runner id '%s' to integer", runnerIDString)
		return projectID, runnerID, err
	}
	return projectID, runnerID, nil

}

func resourceGitlabProjectRunnerEnablementDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	projectID, runnerID, err := projectAndRunnerFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] Delete gitlab project runner %s/%v", projectID, runnerID)

	_, err = client.Runners.DisableProjectRunner(projectID, runnerID)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
