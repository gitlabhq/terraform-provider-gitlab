package gitlab

import (
	"errors"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabPipelineSchedule() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabPipelineScheduleCreate,
		Read:   resourceGitlabPipelineScheduleRead,
		Update: resourceGitlabPipelineScheduleUpdate,
		Delete: resourceGitlabPipelineScheduleDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ref": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cron_timezone": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UTC",
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGitlabPipelineScheduleCreate(d *schema.ResourceData, meta interface{}) error {
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

	PipelineSchedule, _, err := client.PipelineSchedules.CreatePipelineSchedule(project, options)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(PipelineSchedule.ID))

	return resourceGitlabPipelineScheduleRead(d, meta)
}

func resourceGitlabPipelineScheduleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	pipelineScheduleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	log.Printf("[DEBUG] read gitlab PipelineSchedule %s/%d", project, pipelineScheduleID)

	pipelineSchedules, _, err := client.PipelineSchedules.ListPipelineSchedules(project, nil)
	if err != nil {
		return err
	}
	found := false
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
	if !found {
		return errors.New(fmt.Sprintf("PipelineSchedule %d no longer exists in gitlab", pipelineScheduleID))
	}

	return nil
}

func resourceGitlabPipelineScheduleUpdate(d *schema.ResourceData, meta interface{}) error {
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
		return fmt.Errorf("%s cannot be converted to int", d.Id())
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

	_, _, err = client.PipelineSchedules.EditPipelineSchedule(project, pipelineScheduleID, options)
	if err != nil {
		return err
	}

	return resourceGitlabPipelineScheduleRead(d, meta)
}

func resourceGitlabPipelineScheduleDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab PipelineSchedule %s", d.Id())

	pipelineScheduleID, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	resp, err := client.PipelineSchedules.DeletePipelineSchedule(project, pipelineScheduleID)
	if err != nil {
		return fmt.Errorf("%s failed to delete pipeline schedule: %s", d.Id(), resp.Status)
	}
	return err
}
