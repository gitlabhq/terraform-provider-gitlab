package gitlab

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectVariableCreate,
		Read:   resourceGitlabProjectVariableRead,
		Update: resourceGitlabProjectVariableUpdate,
		Delete: resourceGitlabProjectVariableDelete,

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				Required: true,
			},
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGitlabProjectVariableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	options := &gitlab.CreateBuildVariableOptions{
		Key:       gitlab.String(key),
		Value:     gitlab.String(d.Get("value").(string)),
		Protected: gitlab.Bool(d.Get("protected").(bool)),
	}
	log.Printf("[DEBUG] create gitlab project variable %s/%s", project, key)

	_, _, err := client.BuildVariables.CreateBuildVariable(project, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] read gitlab project variable %s/%s", project, key)

	v, _, err := client.BuildVariables.GetBuildVariable(project, key)
	if err != nil {
		return err
	}

	d.Set("value", v.Value)
	d.Set("protected", v.Protected)
	return nil
}

func resourceGitlabProjectVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	options := &gitlab.UpdateBuildVariableOptions{
		Key:       gitlab.String(d.Get("key").(string)),
		Value:     gitlab.String(d.Get("value").(string)),
		Protected: gitlab.Bool(d.Get("protected").(bool)),
	}
	log.Printf("[DEBUG] update gitlab project variable %s/%s", project, key)

	_, _, err := client.BuildVariables.UpdateBuildVariable(project, key, options)
	if err != nil {
		return err
	}

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Delete gitlab project variable %s/%s", project, key)

	_, err := client.BuildVariables.RemoveBuildVariable(project, key)
	return err
}
