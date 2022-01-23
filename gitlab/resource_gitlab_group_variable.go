package gitlab

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabGroupVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitlabGroupVariableCreate,
		ReadContext:   resourceGitlabGroupVariableRead,
		UpdateContext: resourceGitlabGroupVariableUpdate,
		DeleteContext: resourceGitlabGroupVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceGitlabGroupVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, _, err := client.GroupVariables.CreateVariable(group, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(buildTwoPartID(&group, &key))

	return resourceGitlabGroupVariableRead(ctx, d, meta)
}

func resourceGitlabGroupVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group, key, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read gitlab group variable %s/%s", group, key)

	v, _, err := client.GroupVariables.GetVariable(group, key, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab group variable not found %s/%s", group, key)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	d.Set("key", v.Key)
	d.Set("value", v.Value)
	d.Set("variable_type", v.VariableType)
	d.Set("group", group)
	d.Set("protected", v.Protected)
	d.Set("masked", v.Masked)
	return nil
}

func resourceGitlabGroupVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, _, err := client.GroupVariables.UpdateVariable(group, key, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGitlabGroupVariableRead(ctx, d, meta)
}

func resourceGitlabGroupVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Delete gitlab group variable %s/%s", group, key)

	_, err := client.GroupVariables.RemoveVariable(group, key, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
