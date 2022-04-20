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

var _ = registerResource("gitlab_pipeline_trigger", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_pipeline_trigger` + "`" + ` resource allows to manage the lifecycle of a pipeline trigger.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/pipeline_triggers.html)`,

		CreateContext: resourceGitlabPipelineTriggerCreate,
		ReadContext:   resourceGitlabPipelineTriggerRead,
		UpdateContext: resourceGitlabPipelineTriggerUpdate,
		DeleteContext: resourceGitlabPipelineTriggerDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabPipelineTriggerStateImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the trigger to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the pipeline trigger.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"token": {
				Description: "The pipeline trigger token.",
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
			},
		},
	}
})

func resourceGitlabPipelineTriggerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.AddPipelineTriggerOptions{
		Description: gitlab.String(d.Get("description").(string)),
	}

	log.Printf("[DEBUG] create gitlab PipelineTrigger %s", *options.Description)

	PipelineTrigger, _, err := client.PipelineTriggers.AddPipelineTrigger(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(strconv.Itoa(PipelineTrigger.ID))

	return resourceGitlabPipelineTriggerRead(ctx, d, meta)
}

func resourceGitlabPipelineTriggerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	log.Printf("[DEBUG] read gitlab PipelineTrigger %s/%d", project, pipelineTriggerID)

	pipelineTrigger, _, err := client.PipelineTriggers.GetPipelineTrigger(project, pipelineTriggerID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab pipeline trigger not found %s/%d", project, pipelineTriggerID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("description", pipelineTrigger.Description)
	d.Set("token", pipelineTrigger.Token)

	return nil
}

func resourceGitlabPipelineTriggerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.EditPipelineTriggerOptions{
		Description: gitlab.String(d.Get("description").(string)),
	}

	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] update gitlab PipelineTrigger %s", d.Id())

	_, _, err = client.PipelineTriggers.EditPipelineTrigger(project, pipelineTriggerID, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabPipelineTriggerRead(ctx, d, meta)
}

func resourceGitlabPipelineTriggerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab PipelineTrigger %s", d.Id())

	pipelineTriggerID, err := strconv.Atoi(d.Id())

	if err != nil {
		return diag.Errorf("%s cannot be converted to int", d.Id())
	}

	_, err = client.PipelineTriggers.DeletePipelineTrigger(project, pipelineTriggerID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabPipelineTriggerStateImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
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
