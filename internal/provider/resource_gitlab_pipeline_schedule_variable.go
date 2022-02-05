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

var _ = registerResource("gitlab_pipeline_schedule_variable", func() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage variables for pipeline schedules.",

		CreateContext: resourceGitlabPipelineScheduleVariableCreate,
		ReadContext:   resourceGitlabPipelineScheduleVariableRead,
		UpdateContext: resourceGitlabPipelineScheduleVariableUpdate,
		DeleteContext: resourceGitlabPipelineScheduleVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabPipelineScheduleVariableImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The id of the project to add the schedule to.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"pipeline_schedule_id": {
				Description: "The id of the pipeline schedule.",
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
			},
			"key": {
				Description: "Name of the variable.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"value": {
				Description: "Value of the variable.",
				Type:        schema.TypeString,
				Required:    true,
			},
		},
	}
})

func resourceGitlabPipelineScheduleVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	scheduleID := d.Get("pipeline_schedule_id").(int)

	options := &gitlab.CreatePipelineScheduleVariableOptions{
		Key:   gitlab.String(d.Get("key").(string)),
		Value: gitlab.String(d.Get("value").(string)),
	}

	log.Printf("[DEBUG] create gitlab PipelineScheduleVariable %s:%s", *options.Key, *options.Value)

	scheduleVar, _, err := client.PipelineSchedules.CreatePipelineScheduleVariable(project, scheduleID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	id := strconv.Itoa(scheduleID)
	d.SetId(buildTwoPartID(&id, &scheduleVar.Key))

	return resourceGitlabPipelineScheduleVariableRead(ctx, d, meta)
}

func resourceGitlabPipelineScheduleVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	scheduleID := d.Get("pipeline_schedule_id").(int)
	pipelineVariableKey := d.Get("key").(string)

	log.Printf("[DEBUG] read gitlab PipelineSchedule %s/%d", project, scheduleID)

	pipelineSchedule, _, err := client.PipelineSchedules.GetPipelineSchedule(project, scheduleID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	found := false
	for _, pipelineVariable := range pipelineSchedule.Variables {
		if pipelineVariable.Key == pipelineVariableKey {
			d.Set("project", project)
			d.Set("key", pipelineVariable.Key)
			d.Set("value", pipelineVariable.Value)
			d.Set("pipeline_schedule_id", scheduleID)
			found = true
			break
		}
	}
	if !found {
		log.Printf("[DEBUG] pipeline schedule variable not found %s/%d/%s", project, scheduleID, pipelineVariableKey)
		d.SetId("")
	}

	return nil
}

func resourceGitlabPipelineScheduleVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	variableKey := d.Get("key").(string)
	scheduleID := d.Get("pipeline_schedule_id").(int)

	if d.HasChange("value") {
		options := &gitlab.EditPipelineScheduleVariableOptions{
			Value: gitlab.String(d.Get("value").(string)),
		}

		log.Printf("[DEBUG] update gitlab PipelineScheduleVariable %s", d.Id())

		_, _, err := client.PipelineSchedules.EditPipelineScheduleVariable(project, scheduleID, variableKey, options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabPipelineScheduleVariableRead(ctx, d, meta)
}

func resourceGitlabPipelineScheduleVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	variableKey := d.Get("key").(string)
	scheduleID := d.Get("pipeline_schedule_id").(int)

	if _, _, err := client.PipelineSchedules.DeletePipelineScheduleVariable(project, scheduleID, variableKey, gitlab.WithContext(ctx)); err != nil {
		return diag.Errorf("%s failed to delete pipeline schedule variable: %s", d.Id(), err.Error())
	}
	return nil
}

func resourceGitlabPipelineScheduleVariableImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 3 {
		return nil, fmt.Errorf("invalid pipeline schedule variable import format; expected '{project_id}:{pipeline_schedule_id}:{key}'")
	}
	project, pipelineScheduleId, key := s[0], s[1], s[2]
	psid, err := strconv.Atoi(pipelineScheduleId)
	if err != nil {
		return nil, err
	}
	d.SetId(buildTwoPartID(&pipelineScheduleId, &key))
	if err := d.Set("project", project); err != nil {
		return nil, err
	}
	if err := d.Set("pipeline_schedule_id", psid); err != nil {
		return nil, err
	}
	if err := d.Set("key", key); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
