package gitlab

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

func resourceGitlabPipelineSchedule() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage pipeline schedules.\n" +
			"For further information on clusters, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ce/user/project/pipelines/schedules.html).",

		CreateContext: resourceGitlabPipelineScheduleCreate,
		ReadContext:   resourceGitlabPipelineScheduleRead,
		UpdateContext: resourceGitlabPipelineScheduleUpdate,
		DeleteContext: resourceGitlabPipelineScheduleDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabPipelineScheduleStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the schedule to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"description": {
				Description: "The description of the pipeline schedule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"ref": {
				Description: "The branch/tag name to be triggered.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cron": {
				Description: "The cron (e.g. `0 1 * * *`).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cron_timezone": {
				Description: "The timezone.",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "UTC",
			},
			"active": {
				Description: "The activation of pipeline schedule. If false is set, the pipeline schedule will deactivated initially.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
			},
		},
	}
}

func resourceGitlabPipelineScheduleCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreatePipelineScheduleOptions{
		Description:  gitlab.String(d.Get("description").(string)),
		Ref:          gitlab.String(d.Get("ref").(string)),
		Cron:         gitlab.String(d.Get("cron").(string)),
		CronTimezone: gitlab.String(d.Get("cron_timezone").(string)),
		Active:       gitlab.Bool(d.Get("active").(bool)),
	}

	log.Printf("[DEBUG] create gitlab PipelineSchedule %s", *options.Description)

	PipelineSchedule, _, err := client.PipelineSchedules.CreatePipelineSchedule(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(PipelineSchedule.ID))

	return resourceGitlabPipelineScheduleRead(ctx, d, meta)
}

func resourceGitlabPipelineScheduleRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	pipelineScheduleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	log.Printf("[DEBUG] read gitlab PipelineSchedule %s/%d", project, pipelineScheduleID)

	opt := &gitlab.ListPipelineSchedulesOptions{
		Page:    1,
		PerPage: 20,
	}

	found := false
	for {
		pipelineSchedules, resp, err := client.PipelineSchedules.ListPipelineSchedules(project, opt, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		for _, pipelineSchedule := range pipelineSchedules {
			if pipelineSchedule.ID == pipelineScheduleID {
				d.Set("description", pipelineSchedule.Description)
				d.Set("ref", pipelineSchedule.Ref)
				d.Set("cron", pipelineSchedule.Cron)
				d.Set("cron_timezone", pipelineSchedule.CronTimezone)
				d.Set("active", pipelineSchedule.Active)
				found = true
				break
			}
		}

		if found || resp.CurrentPage >= resp.TotalPages {
			break
		}

		opt.Page = resp.NextPage
	}
	if !found {
		log.Printf("[DEBUG] PipelineSchedule %d no longer exists in gitlab", pipelineScheduleID)
		d.SetId("")
		return nil
	}

	return nil
}

func resourceGitlabPipelineScheduleUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.EditPipelineScheduleOptions{
		Description:  gitlab.String(d.Get("description").(string)),
		Ref:          gitlab.String(d.Get("ref").(string)),
		Cron:         gitlab.String(d.Get("cron").(string)),
		CronTimezone: gitlab.String(d.Get("cron_timezone").(string)),
		Active:       gitlab.Bool(d.Get("active").(bool)),
	}

	pipelineScheduleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	if d.HasChange("ref") {
		options.Ref = gitlab.String(d.Get("ref").(string))
	}

	if d.HasChange("cron") {
		options.Cron = gitlab.String(d.Get("cron").(string))
	}

	if d.HasChange("cron_timezone") {
		options.CronTimezone = gitlab.String(d.Get("cron_timezone").(string))
	}

	if d.HasChange("active") {
		options.Active = gitlab.Bool(d.Get("active").(bool))
	}

	log.Printf("[DEBUG] update gitlab PipelineSchedule %s", d.Id())

	_, _, err = client.PipelineSchedules.EditPipelineSchedule(project, pipelineScheduleID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabPipelineScheduleRead(ctx, d, meta)
}

func resourceGitlabPipelineScheduleDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab PipelineSchedule %s", d.Id())

	pipelineScheduleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	if _, err = client.PipelineSchedules.DeletePipelineSchedule(project, pipelineScheduleID, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("failed to delete pipeline schedule %q: %v", d.Id(), err)
	}
	return nil
}

func resourceGitlabPipelineScheduleStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Pipeline Schedule import format; expected '{project_id}:{pipeline_schedule_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
