package gitlab

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectVariableCreate,
		Read:   resourceGitlabProjectVariableRead,
		Update: resourceGitlabProjectVariableUpdate,
		Delete: resourceGitlabProjectVariableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"project": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"key": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: StringIsGitlabVariableName(),
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
	value := d.Get("value").(string)
	protected := d.Get("protected").(bool)

	options := gitlab.CreateVariableOptions{
		Key:              &key,
		Value:            &value,
		Protected:        &protected,
		EnvironmentScope: nil,
	}
	log.Printf("[DEBUG] create gitlab project variable %s/%s", project, key)

	_, _, err := client.ProjectVariables.CreateVariable(project, &options)
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(&project, &key))

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project, key, err := parseTwoPartID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab project variable %s/%s", project, key)

	v, _, err := client.ProjectVariables.GetVariable(project, key)
	if err != nil {
		return err
	}

	d.Set("key", v.Key)
	d.Set("value", v.Value)
	d.Set("project", project)
	d.Set("protected", v.Protected)
	return nil
}

func resourceGitlabProjectVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	protected := d.Get("protected").(bool)

	options := &gitlab.UpdateVariableOptions{
		Value:            &value,
		Protected:        &protected,
		EnvironmentScope: nil,
	}
	log.Printf("[DEBUG] update gitlab project variable %s/%s", project, key)

	_, _, err := client.ProjectVariables.UpdateVariable(project, key, options)
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

	_, err := client.ProjectVariables.RemoveVariable(project, key)
	return err
}
