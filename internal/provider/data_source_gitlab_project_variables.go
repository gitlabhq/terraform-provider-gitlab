package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_variables", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_variables`" + ` data source allows to retrieve all project-level CI/CD variables.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/project_level_variables.html)`,

		ReadContext: dataSourceGitlabProjectVariablesRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project.",
				Type:        schema.TypeString,
				ForceNew:    true,
				Required:    true,
			},
			"environment_scope": {
				Description: "The environment scope of the variable. Defaults to all environment (`*`).",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "*",
			},
			"variables": {
				Description: "The list of variables returned by the search",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabProjectVariableGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabProjectVariablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	environmentScope := d.Get("environment_scope").(string)

	options := &gitlab.ListProjectVariablesOptions{
		Page:    1,
		PerPage: 20,
	}

	var variables []*gitlab.ProjectVariable
	for options.Page != 0 {
		paginatedVariables, resp, err := client.ProjectVariables.ListVariables(project, options, gitlab.WithContext(ctx), withEnvironmentScopeFilter(ctx, environmentScope))
		if err != nil {
			return diag.FromErr(err)
		}

		variables = append(variables, paginatedVariables...)
		options.Page = resp.NextPage
	}

	d.SetId(fmt.Sprintf("%s:%s", project, environmentScope))
	if err := d.Set("variables", flattenGitlabProjectVariables(project, variables)); err != nil {
		return diag.Errorf("failed to set variables to state: %v", err)
	}
	return nil
}

func flattenGitlabProjectVariables(project string, variables []*gitlab.ProjectVariable) (values []map[string]interface{}) {
	for _, variable := range variables {
		values = append(values, gitlabProjectVariableToStateMap(project, variable))
	}
	return values
}
