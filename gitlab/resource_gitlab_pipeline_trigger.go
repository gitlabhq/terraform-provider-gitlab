package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabPipelineTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabPipelineTriggerCreate,
		Read:   resourceGitlabPipelineTriggerRead,
		Update: resourceGitlabPipelineTriggerUpdate,
		Delete: resourceGitlabPipelineTriggerDelete,
		Importer: &schema.ResourceImporter{
			State: resourceGitlabPipelineTriggerStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"token": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceGitlabPipelineTriggerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddPipelineTriggerOptions{
		Description: gitlab.String(d.Get("description").(string)),
	}

	log.Printf("[DEBUG] create gitlab PipelineTrigger %s", *options.Description)

	PipelineTrigger, _, err := client.PipelineTriggers.AddPipelineTrigger(project, options)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(PipelineTrigger.ID))

	return resourceGitlabPipelineTriggerRead(d, meta)
}

func resourceGitlabPipelineTriggerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	log.Printf("[DEBUG] read gitlab PipelineTrigger %s/%d", project, pipelineTriggerID)

	pipelineTrigger, _, err := client.PipelineTriggers.GetPipelineTrigger(project, pipelineTriggerID)
	if err != nil {
		return err
	}

	d.Set("description", pipelineTrigger.Description)
	d.Set("token", pipelineTrigger.Token)

	return nil
}

func resourceGitlabPipelineTriggerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.EditPipelineTriggerOptions{
		Description: gitlab.String(d.Get("description").(string)),
	}

	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] update gitlab PipelineTrigger %s", d.Id())

	_, _, err = client.PipelineTriggers.EditPipelineTrigger(project, pipelineTriggerID, options)
	if err != nil {
		return err
	}

	return resourceGitlabPipelineTriggerRead(d, meta)
}

func resourceGitlabPipelineTriggerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab PipelineTrigger %s", d.Id())

	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return fmt.Errorf("%s cannot be converted to int", d.Id())
	}

	_, err = client.PipelineTriggers.DeletePipelineTrigger(project, pipelineTriggerID)
	return err
}

func resourceGitlabPipelineTriggerStateImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	s := strings.Split(d.Id(), ":")
	if len(s) != 2 {
		d.SetId("")
		return nil, fmt.Errorf("Invalid Pipeline Trigger import format; expected '{project_id}:{pipeline_trigger_id}'")
	}
	project, id := s[0], s[1]

	d.SetId(id)
	d.Set("project", project)

	return []*schema.ResourceData{d}, nil
}
