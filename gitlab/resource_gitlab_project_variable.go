package gitlab

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
			"environment_scope": {
				Type:     schema.TypeString,
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
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)
	environmentScope := d.Get("environment_scope").(string)

	options := gitlab.CreateProjectVariableOptions{
		Key:              &key,
		Value:            &value,
		VariableType:     variableType,
		Protected:        &protected,
		Masked:           &masked,
		EnvironmentScope: &environmentScope,
	}
	log.Printf("[DEBUG] create gitlab project variable %s/%s", project, key)

	_, _, err := client.ProjectVariables.CreateVariable(project, &options)
	if err != nil {
		return augmentProjectVariableClientError(d, err)
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
		return augmentProjectVariableClientError(d, err)
	}

	d.Set("key", v.Key)
	d.Set("value", v.Value)
	d.Set("variable_type", v.VariableType)
	d.Set("project", project)
	d.Set("protected", v.Protected)
	d.Set("masked", v.Masked)
	//For now I'm ignoring environment_scope when reading back data. (this can cause configuration drift so it is bad).
	//However I'm unable to stop terraform from gratuitously updating this to values that are unacceptable by Gitlab)
	//I don't have an enterprise license to properly test this either.
	d.Set("environment_scope", v.EnvironmentScope)
	return nil
}

func resourceGitlabProjectVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	variableType := stringToVariableType(d.Get("variable_type").(string))
	protected := d.Get("protected").(bool)
	masked := d.Get("masked").(bool)
	environmentScope := d.Get("environment_scope").(string)

	options := &gitlab.UpdateProjectVariableOptions{
		Value:            &value,
		VariableType:     variableType,
		Protected:        &protected,
		Masked:           &masked,
		EnvironmentScope: &environmentScope,
	}
	log.Printf("[DEBUG] update gitlab project variable %s/%s", project, key)

	_, _, err := client.ProjectVariables.UpdateVariable(project, key, options)
	if err != nil {
		return augmentProjectVariableClientError(d, err)
	}

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	log.Printf("[DEBUG] Delete gitlab project variable %s/%s", project, key)

	_, err := client.ProjectVariables.RemoveVariable(project, key)
	return augmentProjectVariableClientError(d, err)
}

func augmentProjectVariableClientError(d *schema.ResourceData, err error) error {
	// Masked values will commonly error due to their strict requirements, and the error message from the GitLab API is not very informative,
	// so we return a custom error message in this case.
	if d.Get("masked").(bool) && isInvalidValueError(err) {
		log.Printf("[ERROR] %v", err)
		return errors.New("Invalid value for a masked variable. Check the masked variable requirements: https://docs.gitlab.com/ee/ci/variables/#masked-variable-requirements")
	}

	return err
}

func isInvalidValueError(err error) bool {
	var httpErr *gitlab.ErrorResponse
	return errors.As(err, &httpErr) &&
		httpErr.Response.StatusCode == http.StatusBadRequest &&
		strings.Contains(httpErr.Message, "value") &&
		strings.Contains(httpErr.Message, "invalid")
}
