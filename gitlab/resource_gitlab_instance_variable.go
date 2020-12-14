package gitlab

import (
	"log"
	"net/http"

	gitlab "github.com/Fourcast/go-gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceGitlabInstanceVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabInstanceVariableCreate,
		Read:   resourceGitlabInstanceVariableRead,
		Update: resourceGitlabInstanceVariableUpdate,
		Delete: resourceGitlabInstanceVariableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"key": {
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: StringIsGitlabVariableName,
			},
			"value": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"variable_type": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "env_var",
				ValidateFunc: StringIsGitlabVariableType,
			},
			"protected": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"masked": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func resourceGitlabInstanceVariableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	key := d.Get("key").(string)
	value := d.Get("value").(string)
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)

	options := gitlab.CreateInstanceVariableOptions{
		Key:          &key,
		Value:        &value,
		VariableType: variableType,
		Protected:    &protected,
		Masked:       &masked,
	}
	log.Printf("[DEBUG] create gitlab instance level CI variable %s", key)

	_, _, err := client.InstanceVariables.CreateVariable(&options)
	if err != nil {
		return err
	}

	d.SetId(key)

	return resourceGitlabInstanceVariableRead(d, meta)
}

func resourceGitlabInstanceVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	key := d.Id()

	log.Printf("[DEBUG] read gitlab instance level CI variable %s", key)

	v, resp, err := client.InstanceVariables.GetVariable(key)
	if err != nil {
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("[DEBUG] gitlab instance level CI variable for %s not found so removing from state", d.Id())
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("key", v.Key)
	d.Set("value", v.Value)
	d.Set("variable_type", v.VariableType)
	d.Set("protected", v.Protected)
	d.Set("masked", v.Masked)
	return nil
}

func resourceGitlabInstanceVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	key := d.Get("key").(string)
	value := d.Get("value").(string)
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)

	options := &gitlab.UpdateInstanceVariableOptions{
		Value:        &value,
		Protected:    &protected,
		VariableType: variableType,
		Masked:       &masked,
	}
	log.Printf("[DEBUG] update gitlab instance level CI variable %s", key)

	_, _, err := client.InstanceVariables.UpdateVariable(key, options)
	if err != nil {
		return err
	}
	return resourceGitlabInstanceVariableRead(d, meta)
}

func resourceGitlabInstanceVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Delete gitlab instance level CI variable %s", key)

	_, err := client.InstanceVariables.RemoveVariable(key)
	return err
}
