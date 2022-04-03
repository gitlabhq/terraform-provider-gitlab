package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_variable", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_variable`" + ` data source allows to retrieve details about a project-level CI/CD variable.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/project_level_variables.html)`,

		ReadContext: dataSourceGitlabProjectVariableRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabProjectVariableGetSchema(), []string{"project", "key"}, []string{"environment_scope"}),
	}
})

func dataSourceGitlabProjectVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	key := d.Get("key").(string)
	environmentScope := d.Get("environment_scope").(string)

	variable, _, err := client.ProjectVariables.GetVariable(project, key, nil, gitlab.WithContext(ctx), withEnvironmentScopeFilter(ctx, environmentScope))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", project, key, environmentScope))
	stateMap := gitlabProjectVariableToStateMap(project, variable)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
