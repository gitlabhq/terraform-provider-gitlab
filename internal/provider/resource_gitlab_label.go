package gitlab

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabLabel() *schema.Resource {
	// lintignore: XR002 // TODO: Resolve this tfproviderlint issue
	return &schema.Resource{
		Description: "This resource allows you to create and manage labels for your GitLab projects.\n" +
			"For further information on labels, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ee/user/project/labels.html#project-labels).",

		CreateContext: resourceGitlabLabelCreate,
		ReadContext:   resourceGitlabLabelRead,
		UpdateContext: resourceGitlabLabelUpdate,
		DeleteContext: resourceGitlabLabelDelete,

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
}

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

	page := 1
	labelsLen := 0
	for page == 1 || labelsLen != 0 {
		labels, _, err := client.Labels.ListLabels(project, &gitlab.ListLabelsOptions{ListOptions: gitlab.ListOptions{Page: page}}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		for _, label := range labels {
			if label.Name == labelName {
				d.Set("description", label.Description)
				d.Set("color", label.Color)
				d.Set("name", label.Name)
				return nil
			}
		}
		labelsLen = len(labels)
		page = page + 1
	}

	log.Printf("[DEBUG] failed to read gitlab label %s/%s", project, labelName)
	d.SetId("")
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
