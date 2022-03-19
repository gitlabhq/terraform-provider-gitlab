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

var _ = registerResource("gitlab_label", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_label` + "`" + ` resource allows to manage the lifecycle of a project label.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/labels.html#project-labels)`,

		CreateContext: resourceGitlabLabelCreate,
		ReadContext:   resourceGitlabLabelRead,
		UpdateContext: resourceGitlabLabelUpdate,
		DeleteContext: resourceGitlabLabelDelete,
		// FIXME: this importer sucks a little, but we can't use a passthrough importer, because
		//        the resource id is flawed and we don't want to break backwards-compatibility.
		//        We cannot have the same label in two different projects as of now, ...
		Importer: &schema.ResourceImporter{
			StateContext: resourceGitlabLabelImporter,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project to add the label to.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the label.",
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
			},
			"color": {
				Description: "The color of the label given in 6-digit hex notation with leading '#' sign (e.g. #FFAABB) or one of the [CSS color names](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value#Color_keywords).",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the label.",
				Type:        schema.TypeString,
				Optional:    true,
			},
		},
	}
})

func resourceGitlabLabelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.CreateLabelOptions{
		Name:  gitlab.String(d.Get("name").(string)),
		Color: gitlab.String(d.Get("color").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab label %s", *options.Name)

	label, _, err := client.Labels.CreateLabel(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(label.Name)

	return resourceGitlabLabelRead(ctx, d, meta)
}

func resourceGitlabLabelRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	labelName := d.Id()
	log.Printf("[DEBUG] read gitlab label %s/%s", project, labelName)

	label, _, err := client.Labels.GetLabel(project, labelName, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] failed to read gitlab label %s/%s", project, labelName)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("description", label.Description)
	d.Set("color", label.Color)
	d.Set("name", label.Name)
	return nil
}

func resourceGitlabLabelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	options := &gitlab.UpdateLabelOptions{
		Name:  gitlab.String(d.Get("name").(string)),
		Color: gitlab.String(d.Get("color").(string)),
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] update gitlab label %s", d.Id())

	_, _, err := client.Labels.UpdateLabel(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceGitlabLabelRead(ctx, d, meta)
}

func resourceGitlabLabelDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab label %s", d.Id())
	options := &gitlab.DeleteLabelOptions{
		Name: gitlab.String(d.Id()),
	}

	_, err := client.Labels.DeleteLabel(project, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabLabelImporter(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*gitlab.Client)
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid label id (should be <project ID>.<label name>): %s", d.Id())
	}

	d.SetId(parts[1])
	project, _, err := client.Projects.GetProject(parts[0], nil, gitlab.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	if err := d.Set("project", strconv.Itoa(project.ID)); err != nil {
		return nil, err
	}

	diagnostic := resourceGitlabLabelRead(ctx, d, meta)
	if diagnostic.HasError() {
		return nil, fmt.Errorf("failed to read project label %s: %s", d.Id(), diagnostic[0].Summary)
	}

	return []*schema.ResourceData{d}, nil
}
