package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_variable", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_variable`" + ` data source allows to retrieve details about a group-level CI/CD variable.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/group_level_variables.html)`,

		ReadContext: dataSourceGitlabGroupVariableRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabGroupVariableGetSchema(), []string{"group", "key"}, []string{"environment_scope"}),
	}
})

func dataSourceGitlabGroupVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	key := d.Get("key").(string)
	environmentScope := d.Get("environment_scope").(string)

	variable, _, err := client.GroupVariables.GetVariable(group, key, nil, gitlab.WithContext(ctx), withEnvironmentScopeFilter(ctx, environmentScope))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%s:%s", group, key, environmentScope))
	stateMap := gitlabGroupVariableToStateMap(group, variable)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
