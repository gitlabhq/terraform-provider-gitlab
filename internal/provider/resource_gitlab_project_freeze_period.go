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

var _ = registerResource("gitlab_project_freeze_period", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_project_freeze_period` + "`" + ` resource allows to manage the lifecycle of a freeze period for a project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/freeze_periods.html)`,

		CreateContext: resourceGitlabProjectFreezePeriodCreate,
		ReadContext:   resourceGitlabProjectFreezePeriodRead,
		UpdateContext: resourceGitlabProjectFreezePeriodUpdate,
		DeleteContext: resourceGitlabProjectFreezePeriodDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description: "The id of the project to add the schedule to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"freeze_start": {
				Description: "Start of the Freeze Period in cron format (e.g. `0 1 * * *`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"freeze_end": {
				Description: "End of the Freeze Period in cron format (e.g. `0 2 * * *`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cron_timezone": {
				Description: "The timezone.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "UTC",
			},
		},
	}
})

func resourceGitlabProjectFreezePeriodCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	projectID := d.Get("project_id").(string)

	options := gitlab.CreateFreezePeriodOptions{
		FreezeStart:  gitlab.String(d.Get("freeze_start").(string)),
		FreezeEnd:    gitlab.String(d.Get("freeze_end").(string)),
		CronTimezone: gitlab.String(d.Get("cron_timezone").(string)),
	}

	log.Printf("[DEBUG] Project %s create gitlab project-level freeze period %+v", projectID, options)

	client := meta.(*gitlab.Client)
	FreezePeriod, _, err := client.FreezePeriods.CreateFreezePeriodOptions(projectID, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	FreezePeriodIDString := fmt.Sprintf("%d", FreezePeriod.ID)
	d.SetId(buildTwoPartID(&projectID, &FreezePeriodIDString))

	return resourceGitlabProjectFreezePeriodRead(ctx, d, meta)
}

func resourceGitlabProjectFreezePeriodRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab FreezePeriod %s/%d", projectID, freezePeriodID)

	freezePeriod, _, err := client.FreezePeriods.GetFreezePeriod(projectID, freezePeriodID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] project freeze period for %s not found so removing it from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("freeze_start", freezePeriod.FreezeStart)
	d.Set("freeze_end", freezePeriod.FreezeEnd)
	d.Set("cron_timezone", freezePeriod.CronTimezone)
	d.Set("project_id", projectID)

	return nil
}

func resourceGitlabProjectFreezePeriodUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	options := &gitlab.UpdateFreezePeriodOptions{}

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	if d.HasChange("freeze_start") {
		options.FreezeStart = gitlab.String(d.Get("freeze_start").(string))
	}

	if d.HasChange("freeze_end") {
		options.FreezeEnd = gitlab.String(d.Get("freeze_end").(string))
	}

	if d.HasChange("cron_timezone") {
		options.CronTimezone = gitlab.String(d.Get("cron_timezone").(string))
	}

	log.Printf("[DEBUG] update gitlab FreezePeriod %s", d.Id())

	_, _, err = client.FreezePeriods.UpdateFreezePeriodOptions(projectID, freezePeriodID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabProjectFreezePeriodRead(ctx, d, meta)
}

func resourceGitlabProjectFreezePeriodDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	projectID, freezePeriodID, err := projectIDAndFreezePeriodIDFromID(d.Id())
	log.Printf("[DEBUG] Delete gitlab FreezePeriod %s", d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	if _, err = client.FreezePeriods.DeleteFreezePeriod(projectID, freezePeriodID, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("failed to delete pipeline schedule %q: %v", d.Id(), err)
	}

	return nil
}

func projectIDAndFreezePeriodIDFromID(id string) (string, int, error) {
	project, freezePeriodIDString, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	freezePeriodID, err := strconv.Atoi(freezePeriodIDString)
	if err != nil {
		return "", 0, fmt.Errorf("failed to get freezePeriodId: %v", err)
	}

	return project, freezePeriodID, nil
}
