package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

var gitlabVariableTypeValues = []string{"env_var", "file"}

func gitlabProjectVariableGetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project": {
			Description: "The name or id of the project.",
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
			Description:      fmt.Sprintf("The type of a variable. Valid values are: %s. Default is `env_var`.", renderValueListForDocs(gitlabVariableTypeValues)),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "env_var",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(gitlabVariableTypeValues, false)),
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
			Description: "The environment scope of the variable. Defaults to all environment (`*`). Note that in Community Editions of Gitlab, values other than `*` will cause inconsistent plans.",
			Type:        schema.TypeString,
			Optional:    true,
			Default:     "*",
			// Versions of GitLab prior to 13.4 cannot update environment_scope.
			ForceNew: true,
		},
	}
}

func gitlabProjectVariableToStateMap(project string, variable *gitlab.ProjectVariable) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["key"] = variable.Key
	stateMap["value"] = variable.Value
	stateMap["variable_type"] = variable.VariableType
	stateMap["protected"] = variable.Protected
	stateMap["masked"] = variable.Masked
	stateMap["environment_scope"] = variable.EnvironmentScope
	return stateMap
}
