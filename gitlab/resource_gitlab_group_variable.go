package gitlab

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// modifyRequestAddEnvironmentFilter returns a RequestOptionFunc function that
// can be passed to the go-gitlab library calls to add the environment scope to
// requests to lookup, modification, and deletion requests. Since gitlab 13.11,
// an environment variable key is no longer unique and is composit-keyed with
// the scope.
// See https://docs.gitlab.com/ee/ci/variables/#add-a-cicd-variable-to-a-group
func modifyRequestAddEnvironmentFilter(scope string) gitlab.RequestOptionFunc {
	return func(r *retryablehttp.Request) error {
		queryParams := r.URL.Query()
		queryParams.Add("filter[environment_scope]", scope)
		r.URL.RawQuery = queryParams.Encode()
		return nil
	}
}

func resourceGitlabGroupVariable() *schema.Resource {
	return &schema.Resource{
		Description: "This resource allows you to create and manage CI/CD variables for your GitLab groups.\n" +
			"For further information on variables, consult the [gitlab\n" +
			"documentation](https://docs.gitlab.com/ce/ci/variables/README.html#variables).",

		CreateContext: resourceGitlabGroupVariableCreate,
		ReadContext:   resourceGitlabGroupVariableRead,
		UpdateContext: resourceGitlabGroupVariableUpdate,
		DeleteContext: resourceGitlabGroupVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The name or id of the group.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"key": {
				Description:  "The name of the variable.",
				Type:         schema.TypeString,
				ForceNew:     true,
				Required:     true,
				ValidateFunc: StringIsGitlabVariableName,
			},
			"value": {
				Description: "The value of the variable.",
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
			},
			"variable_type": {
				Description:  "The type of a variable. Available types are: env_var (default) and file.",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "env_var",
				ValidateFunc: StringIsGitlabVariableType,
			},
			"protected": {
				Description: "If set to `true`, the variable will be passed only to pipelines running on protected branches and tags. Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"masked": {
				Description: "If set to `true`, the value of the variable will be hidden in job logs. The value must meet the [masking requirements](https://docs.gitlab.com/ee/ci/variables/#masked-variables). Defaults to `false`.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"environment_scope": {
				Description: "The environment scope of the variable. Defaults to all environment (`*`). Note that in Community Editions of Gitlab, values other than `*` will cause inconsistent plans. See https://docs.gitlab.com/ee/ci/variables/#add-a-cicd-variable-to-a-group",
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Default:     "*",
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
	environmentScope := d.Get("environment_scope").(string)

	options := gitlab.CreateGroupVariableOptions{
		Key:              &key,
		Value:            &value,
		VariableType:     variableType,
		Protected:        &protected,
		Masked:           &masked,
		EnvironmentScope: &environmentScope,
	}
	log.Printf("[DEBUG] create gitlab group variable %s/%s", group, key)

	_, _, err := client.GroupVariables.CreateVariable(group, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	keyScope := fmt.Sprintf("%s:%s", key, environmentScope)

	d.SetId(buildTwoPartID(&group, &keyScope))

	return resourceGitlabGroupVariableRead(ctx, d, meta)
}

func resourceGitlabGroupVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group, key, err := parseTwoPartID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	keyScope := strings.SplitN(key, ":", 2)
	scope := "*"
	if len(keyScope) == 2 {
		key = keyScope[0]
		scope = keyScope[1]
	}

	log.Printf("[DEBUG] read gitlab group variable %s/%s/%s", group, key, scope)

	v, _, err := client.GroupVariables.GetVariable(
		group,
		key,
		gitlab.WithContext(ctx),
		modifyRequestAddEnvironmentFilter(scope),
	)
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
	d.Set("environment_scope", v.EnvironmentScope)
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
	environmentScope := d.Get("environment_scope").(string)

	options := &gitlab.UpdateGroupVariableOptions{
		Value:            &value,
		Protected:        &protected,
		VariableType:     variableType,
		Masked:           &masked,
		EnvironmentScope: &environmentScope,
	}
	log.Printf("[DEBUG] update gitlab group variable %s/%s/%s", group, key, environmentScope)

	_, _, err := client.GroupVariables.UpdateVariable(
		group,
		key,
		options,
		gitlab.WithContext(ctx),
		modifyRequestAddEnvironmentFilter(environmentScope),
	)
	if err != nil {
		return diag.FromErr(err)
	}
	return resourceGitlabGroupVariableRead(ctx, d, meta)
}

func resourceGitlabGroupVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	key := d.Get("key").(string)
	environmentScope := d.Get("environment_scope").(string)
	log.Printf("[DEBUG] Delete gitlab group variable %s/%s/%s", group, key, environmentScope)

	_, err := client.GroupVariables.RemoveVariable(
		group,
		key,
		gitlab.WithContext(ctx),
		modifyRequestAddEnvironmentFilter(environmentScope),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}
