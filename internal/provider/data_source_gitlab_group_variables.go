package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_variables", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_variables`" + ` data source allows to retrieve all group-level CI/CD variables.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/group_level_variables.html)`,

		ReadContext: dataSourceGitlabGroupVariablesRead,
		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The name or id of the group.",
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
					Schema: datasourceSchemaFromResourceSchema(gitlabGroupVariableGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabGroupVariablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	environmentScope := d.Get("environment_scope").(string)

	options := &gitlab.ListGroupVariablesOptions{
		Page:    1,
		PerPage: 20,
	}

	var variables []*gitlab.GroupVariable
	for options.Page != 0 {
		paginatedVariables, resp, err := client.GroupVariables.ListVariables(group, options, gitlab.WithContext(ctx), withEnvironmentScopeFilter(ctx, environmentScope))
		if err != nil {
			return diag.FromErr(err)
		}

		variables = append(variables, paginatedVariables...)
		options.Page = resp.NextPage
	}

	d.SetId(fmt.Sprintf("%s:%s", group, environmentScope))
	if err := d.Set("variables", flattenGitlabGroupVariables(group, variables)); err != nil {
		return diag.Errorf("failed to set variables to state: %v", err)
	}
	return nil
}

func flattenGitlabGroupVariables(group string, variables []*gitlab.GroupVariable) (values []map[string]interface{}) {
	for _, variable := range variables {
		values = append(values, gitlabGroupVariableToStateMap(group, variable))
	}
	return values
}
