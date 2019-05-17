package gitlab

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabLabelCreate,
		Read:   resourceGitlabLabelRead,
		Update: resourceGitlabLabelUpdate,
		Delete: resourceGitlabLabelDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"color": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitlabLabelCreate(d *schema.ResourceData, meta interface{}) error {
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

	label, _, err := client.Labels.CreateLabel(project, options)
	if err != nil {
		return err
	}

	d.SetId(label.Name)

	return resourceGitlabLabelRead(d, meta)
}

func resourceGitlabLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	labelName := d.Id()
	log.Printf("[DEBUG] read gitlab label %s/%s", project, labelName)

	labels, _, err := client.Labels.ListLabels(project, nil)
	if err != nil {
		return err
	}
	found := false
	for _, label := range labels {
		if label.Name == labelName {
			d.Set("description", label.Description)
			d.Set("color", label.Color)
			d.Set("name", label.Name)
			found = true
			break
		}
	}
	if !found {
		log.Printf("[WARN] removing label %s from state because it no longer exists in gitlab", labelName)
		d.SetId("")
	}

	return nil
}

func resourceGitlabLabelUpdate(d *schema.ResourceData, meta interface{}) error {
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

	_, _, err := client.Labels.UpdateLabel(project, options)
	if err != nil {
		return err
	}

	return resourceGitlabLabelRead(d, meta)
}

func resourceGitlabLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	log.Printf("[DEBUG] Delete gitlab label %s", d.Id())
	options := &gitlab.DeleteLabelOptions{
		Name: gitlab.String(d.Id()),
	}

	_, err := client.Labels.DeleteLabel(project, options)
	return err
}
