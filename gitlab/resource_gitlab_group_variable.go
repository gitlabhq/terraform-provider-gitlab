package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabGroupVariableCreate,
		Read:   resourceGitlabGroupVariableRead,
		Update: resourceGitlabGroupVariableUpdate,
		Delete: resourceGitlabGroupVariableDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
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

func resourceGitlabGroupVariableCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)

	options := gitlab.CreateGroupVariableOptions{
		Key:          &key,
		Value:        &value,
		VariableType: variableType,
		Protected:    &protected,
		Masked:       &masked,
	}
	log.Printf("[DEBUG] create gitlab group variable %s/%s", group, key)

	_, _, err := client.GroupVariables.CreateVariable(group, &options)
	if err != nil {
		return err
	}

	d.SetId(buildTwoPartID(&group, &key))

	return resourceGitlabGroupVariableRead(d, meta)
}

func resourceGitlabGroupVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group, key, err := parseTwoPartID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab group variable %s/%s", group, key)

	v, _, err := client.GroupVariables.GetVariable(group, key)
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group variable not found %s/%s", group, key)
			d.SetId("")
			return nil
		}
		return err
	}

	_ = d.Set("key", v.Key)
	_ = d.Set("value", v.Value)
	_ = d.Set("variable_type", v.VariableType)
	_ = d.Set("group", group)
	_ = d.Set("protected", v.Protected)
	_ = d.Set("masked", v.Masked)
	return nil
}

func resourceGitlabGroupVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)

	options := &gitlab.UpdateGroupVariableOptions{
		Value:        &value,
		Protected:    &protected,
		VariableType: variableType,
		Masked:       &masked,
	}
	log.Printf("[DEBUG] update gitlab group variable %s/%s", group, key)

	_, _, err := client.GroupVariables.UpdateVariable(group, key, options)
	if err != nil {
		return err
	}
	return resourceGitlabGroupVariableRead(d, meta)
}

func resourceGitlabGroupVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Delete gitlab group variable %s/%s", group, key)

	_, err := client.GroupVariables.RemoveVariable(group, key)
	return err
}
