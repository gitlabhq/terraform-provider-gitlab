package gitlab

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabPipelineScheduleVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabPipelineScheduleVariableCreate,
		Read:   resourceGitlabPipelineScheduleVariableRead,
		Update: resourceGitlabPipelineScheduleVariableUpdate,
		Delete: resourceGitlabPipelineScheduleVariableDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"pipeline_schedule_id": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceGitlabPipelineScheduleVariableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	scheduleId := d.Get("pipeline_schedule_id").(int)

	options := &gitlab.CreatePipelineScheduleVariableOptions{
		Key:   gitlab.String(d.Get("key").(string)),
		Value: gitlab.String(d.Get("value").(string)),
	}

	log.Printf("[DEBUG] create gitlab PipelineScheduleVariable %s:%s", *options.Key, *options.Value)

	scheduleVar, _, err := client.PipelineSchedules.CreatePipelineScheduleVariable(project, scheduleId, options)
	if err != nil {
		return err
	}

	id := strconv.Itoa(scheduleId)
	d.SetId(buildTwoPartID(&id, &scheduleVar.Key))

	return resourceGitlabPipelineScheduleVariableRead(d, meta)
}

func resourceGitlabPipelineScheduleVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	scheduleId := d.Get("pipeline_schedule_id").(int)
	pipelineVariableKey := d.Get("key").(string)

	log.Printf("[DEBUG] read gitlab PipelineSchedule %s/%d", project, scheduleId)

	pipelineSchedule, _, err := client.PipelineSchedules.GetPipelineSchedule(project, scheduleId)
	if err != nil {
		return err
	}

	found := false
	for _, pipelineVariable := range pipelineSchedule.Variables {
		if pipelineVariable.Key == pipelineVariableKey {
			d.Set("project", project)
			d.Set("key", pipelineVariable.Key)
			d.Set("value", pipelineVariable.Value)
			d.Set("pipeline_schedule_id", scheduleId)
			found = true
			break
		}
	}
	if !found {
		return errors.New(fmt.Sprintf("PipelineScheduleVariable %s no longer exists", pipelineVariableKey))
	}

	return nil
}

func resourceGitlabPipelineScheduleVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	variableKey := d.Get("key").(string)
	scheduleId := d.Get("pipeline_schedule_id").(int)

	if d.HasChange("value") {
		options := &gitlab.EditPipelineScheduleVariableOptions{
			Value: gitlab.String(d.Get("value").(string)),
		}

		log.Printf("[DEBUG] update gitlab PipelineScheduleVariable %s", d.Id())

		_, _, err := client.PipelineSchedules.EditPipelineScheduleVariable(project, scheduleId, variableKey, options)
		if err != nil {
			return err
		}
	}

	return resourceGitlabPipelineScheduleVariableRead(d, meta)
}

func resourceGitlabPipelineScheduleVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	variableKey := d.Get("key").(string)
	scheduleId := d.Get("pipeline_schedule_id").(int)

	_, resp, err := client.PipelineSchedules.DeletePipelineScheduleVariable(project, scheduleId, variableKey)
	if err != nil {
		return fmt.Errorf("%s failed to delete pipeline schedule variable: %s", d.Id(), resp.Status)
	}
	return err
}
