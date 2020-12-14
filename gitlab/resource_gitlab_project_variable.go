package gitlab

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	gitlab "github.com/Fourcast/go-gitlab"
	retryablehttp "github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Default:  "*",
				// Versions of GitLab prior to 13.4 cannot update environment_scope.
				ForceNew: true,
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

	id := strings.Join([]string{project, key, environmentScope}, ":")

	log.Printf("[DEBUG] create gitlab project variable %q", id)

	_, _, err := client.ProjectVariables.CreateVariable(project, &options)
	if err != nil {
		return augmentProjectVariableClientError(d, err)
	}

	d.SetId(id)

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var (
		project          string
		key              string
		environmentScope string
	)

	// An older version of this resource used the ID format "project:key".
	// For backwards compatibility we still support the old format.
	parts := strings.SplitN(d.Id(), ":", 4)
	switch len(parts) {
	case 2:
		project = parts[0]
		key = parts[1]
		environmentScope = d.Get("environment_scope").(string)
	case 3:
		project = parts[0]
		key = parts[1]
		environmentScope = parts[2]
	default:
		return fmt.Errorf(`Failed to parse project variable ID %q: expected format project:key or project:key:environment_scope`, d.Id())
	}

	log.Printf("[DEBUG] read gitlab project variable %q", d.Id())

	v, err := getProjectVariable(client, project, key, environmentScope)
	if err != nil {
		if errors.Is(err, errProjectVariableNotExist) {
			log.Printf("[DEBUG] read gitlab project variable %q was not found", d.Id())
			d.SetId("")
			return nil
		}
		return augmentProjectVariableClientError(d, err)
	}

	d.Set("key", v.Key)
	d.Set("value", v.Value)
	d.Set("variable_type", v.VariableType)
	d.Set("project", project)
	d.Set("protected", v.Protected)
	d.Set("masked", v.Masked)
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
	log.Printf("[DEBUG] update gitlab project variable %q", d.Id())

	_, _, err := client.ProjectVariables.UpdateVariable(project, key, options, withEnvironmentScopeFilter(environmentScope))
	if err != nil {
		return augmentProjectVariableClientError(d, err)
	}

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	environmentScope := d.Get("environment_scope").(string)
	log.Printf("[DEBUG] Delete gitlab project variable %q", d.Id())

	// Note that the environment_scope filter is added here to support GitLab versions >= 13.4,
	// but it will be ignored in prior versions, causing nondeterministic destroy behavior when
	// destroying or updating scoped variables.
	// ref: https://gitlab.com/gitlab-org/gitlab/-/merge_requests/39209
	_, err := client.ProjectVariables.RemoveVariable(project, key, withEnvironmentScopeFilter(environmentScope))
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

func withEnvironmentScopeFilter(environmentScope string) gitlab.RequestOptionFunc {
	return func(req *retryablehttp.Request) error {
		query, err := url.ParseQuery(req.Request.URL.RawQuery)
		if err != nil {
			return err
		}
		query.Set("filter[environment_scope]", environmentScope)
		req.Request.URL.RawQuery = query.Encode()
		return nil
	}
}

var errProjectVariableNotExist = errors.New("project variable does not exist")

func getProjectVariable(client *gitlab.Client, project interface{}, key, environmentScope string) (*gitlab.ProjectVariable, error) {
	// List and filter variables manually to support GitLab versions < v13.4 (2020-08-22)
	// ref: https://gitlab.com/gitlab-org/gitlab/-/merge_requests/39209

	page := 1

	for {
		projectVariables, resp, err := client.ProjectVariables.ListVariables(project, &gitlab.ListProjectVariablesOptions{Page: page})
		if err != nil {
			return nil, err
		}

		for _, v := range projectVariables {
			if v.Key == key && v.EnvironmentScope == environmentScope {
				return v, nil
			}
		}

		if resp.NextPage == 0 {
			return nil, errProjectVariableNotExist
		}

		page++
	}
}
