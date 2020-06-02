package gitlab

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectVariable() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabProjectVariableCreate,
		Read:   resourceGitlabProjectVariableRead,
		Update: resourceGitlabProjectVariableUpdate,
		Delete: resourceGitlabProjectVariableDelete,
		Exists: resourceGitlabProjectVariableExists,
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

func resourceGitlabProjectVariableExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*gitlab.Client)

	project := d.Get("project")

	log.Printf("[DEBUG] list gitlab project variable %s", project)

	options := gitlab.ListProjectVariablesOptions{Page: 1, PerPage: 9999}
	projectVariables, _, err := client.ProjectVariables.ListVariables(project, &options)
	if err != nil {
		return false, err
	}

	log.Printf("[DEBUG] gitlab project variables: %s", projectVariables)

	key := d.Get("key")
	environmentScope := d.Get("environment_scope")

	for _, projectVariable := range projectVariables {
		if projectVariable.Key == key && projectVariable.EnvironmentScope == environmentScope {
			log.Printf("[DEBUG] Variable matching key and environment scope exists: %s:%s", key, environmentScope)
			return true, nil
		}
	}

	log.Printf("[DEBUG] Variable matching key and environment scope does NOT exist: %s:%s", key, environmentScope)
	return false, nil
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
	log.Printf("[DEBUG] create gitlab project variable %s/%s/%s", project, key, environmentScope)

	_, _, err := client.ProjectVariables.CreateVariable(project, &options)
	if err != nil {
		return err
	}

	d.SetId(buildThreePartID(&project, &key, &environmentScope))

	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	project, key, environmentScope, err := parseThreePartID(d.Id())
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] read gitlab project variable %s/%s", project, key)

	options := gitlab.ListProjectVariablesOptions{Page: 1, PerPage: 9999}
	projectVariables, _, err := client.ProjectVariables.ListVariables(project, &options)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] gitlab project variables: %s", projectVariables)

	for _, projectVariable := range projectVariables {
		if projectVariable.Key == key && projectVariable.EnvironmentScope == environmentScope {
			d.Set("key", projectVariable.Key)
			d.Set("value", projectVariable.Value)
			d.Set("variable_type", projectVariable.VariableType)
			d.Set("project", project)
			d.Set("protected", projectVariable.Protected)
			d.Set("masked", projectVariable.Masked)
			d.Set("environment_scope", projectVariable.EnvironmentScope)
			return nil
		}
	}
	return nil
}

func resourceGitlabProjectVariableUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project")
	key := d.Get("key").(string)
	environmentScopeOfProjectVariableBeingUpdated := d.Get("environment_scope")

	projectVariablesWithSameKeyAsTheOneBeingUpdated, fetchErr := projectVariablesMatchingKey(client, project, key)
	if fetchErr != nil {
		return fetchErr
	}

	deleteErr := deleteProjectVariables(projectVariablesWithSameKeyAsTheOneBeingUpdated, project, client)
	if deleteErr != nil {
		return deleteErr
	}

	for _, projectVariable := range projectVariablesWithSameKeyAsTheOneBeingUpdated {
		environmentScope := projectVariable.EnvironmentScope
		var value string
		var variableType gitlab.VariableTypeValue
		var protected bool
		var masked bool
		if environmentScope == environmentScopeOfProjectVariableBeingUpdated {
			value = d.Get("value").(string)
			variableType = *stringToVariableType(d.Get("variable_type").(string))
			protected = d.Get("protected").(bool)
			masked = d.Get("masked").(bool)
		} else {
			value = projectVariable.Value
			variableType = projectVariable.VariableType
			protected = projectVariable.Protected
			masked = projectVariable.Masked
		}
		options := gitlab.CreateProjectVariableOptions{
			Key:              &key,
			Value:            &value,
			VariableType:     &variableType,
			Protected:        &protected,
			Masked:           &masked,
			EnvironmentScope: &environmentScope,
		}
		log.Printf("[DEBUG] create gitlab project variable %s/%s/%s", project, key, environmentScope)
		_, _, err := client.ProjectVariables.CreateVariable(project, &options)
		if err != nil {
			return err
		}
	}
	return resourceGitlabProjectVariableRead(d, meta)
}

func resourceGitlabProjectVariableDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	project := d.Get("project")
	key := d.Get("key").(string)
	environmentScopeOfProjectVariableBeingDeleted := d.Get("environment_scope")

	projectVariablesWithSameKeyAsTheOneBeingDeleted, fetchErr := projectVariablesMatchingKey(client, project, key)
	if fetchErr != nil {
		return fetchErr
	}

	deleteErr := deleteProjectVariables(projectVariablesWithSameKeyAsTheOneBeingDeleted, project, client)
	if deleteErr != nil {
		return deleteErr
	}

	for _, projectVariable := range projectVariablesWithSameKeyAsTheOneBeingDeleted {
		if projectVariable.EnvironmentScope == environmentScopeOfProjectVariableBeingDeleted {
			continue // Don't re-create project variable for the environment scope that is marked for deletion
		}
		options := gitlab.CreateProjectVariableOptions{
			Key:              &projectVariable.Key,
			Value:            &projectVariable.Value,
			VariableType:     &projectVariable.VariableType,
			Protected:        &projectVariable.Protected,
			Masked:           &projectVariable.Masked,
			EnvironmentScope: &projectVariable.EnvironmentScope,
		}
		log.Printf("[DEBUG] create gitlab project variable %s/%s/%s", project, key, projectVariable.EnvironmentScope)
		_, _, err := client.ProjectVariables.CreateVariable(project, &options)
		if err != nil {
			return err
		}
	}
	return nil
}

func deleteProjectVariables(projectVariables []*gitlab.ProjectVariable, project interface{}, client *gitlab.Client) error {
	for _, projectVariable := range projectVariables {
		log.Printf("[DEBUG] Delete gitlab project variable %s/%s", project, projectVariable.Key)
		_, err := client.ProjectVariables.RemoveVariable(project, projectVariable.Key)
		if err != nil {
			return err
		}
	}
	return nil
}

func projectVariablesMatchingKey(client *gitlab.Client, project interface{}, key string) ([]*gitlab.ProjectVariable, error) {
	options := gitlab.ListProjectVariablesOptions{Page: 1, PerPage: 9999}
	allProjectVariables, _, err := client.ProjectVariables.ListVariables(project, &options)
	if err != nil {
		return nil, err
	}

	var projectVariablesWithSameKey []*gitlab.ProjectVariable

	for _, projectVariable := range allProjectVariables {
		if projectVariable.Key == key {
			projectVariablesWithSameKey = append(projectVariablesWithSameKey, projectVariable)
		}
	}
	return projectVariablesWithSameKey, nil
}
