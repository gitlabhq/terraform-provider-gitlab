package gitlab

import (
	"errors"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupLabel() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupLabelCreate,
		Read:   resourceGitlabGroupLabelRead,
		Update: resourceGitlabGroupLabelUpdate,
		Delete: resourceGitlabGroupLabelDelete,

		Schema: map[string]*schema.Schema{
			"group": {
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

func resourceGitlabGroupLabelCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	options := &gitlab.CreateGroupLabelOptions{
		Name:  gitlab.String(d.Get("name").(string)),
		Color: gitlab.String(d.Get("color").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create gitlab group label %s", *options.Name)

	label, _, err := client.GroupLabels.CreateGroupLabel(group, options)
	if err != nil {
		return err
	}

	d.SetId(label.Name)

	return resourceGitlabGroupLabelRead(d, meta)
}

func resourceGitlabGroupLabelRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	labelName := d.Id()
	log.Printf("[DEBUG] read gitlab group label %s/%s", group, labelName)

	labels, _, err := client.GroupLabels.ListGroupLabels(group, nil)
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
		return errors.New(fmt.Sprintf("label %s does not exist or the user does not have permissions to see it", labelName))
	}

	return nil
}

func resourceGitlabGroupLabelUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	options := &gitlab.UpdateGroupLabelOptions{
		Name:  gitlab.String(d.Get("name").(string)),
		Color: gitlab.String(d.Get("color").(string)),
	}

	if d.HasChange("description") {
		options.Description = gitlab.String(d.Get("description").(string))
	}

	log.Printf("[DEBUG] update gitlab group label %s", d.Id())

	_, _, err := client.GroupLabels.UpdateGroupLabel(group, options)
	if err != nil {
		return err
	}

	return resourceGitlabGroupLabelRead(d, meta)
}

func resourceGitlabGroupLabelDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	log.Printf("[DEBUG] Delete gitlab group label %s", d.Id())
	options := &gitlab.DeleteGroupLabelOptions{
		Name: gitlab.String(d.Id()),
	}

	_, err := client.GroupLabels.DeleteGroupLabel(group, options)
	return err
}
