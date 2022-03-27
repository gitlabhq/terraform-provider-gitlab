package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_instance_variables", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_instance_variables`" + ` data source allows to retrieve all instance-level CI/CD variables.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/instance_level_ci_variables.html)`,

		ReadContext: dataSourceGitlabInstanceVariablesRead,
		Schema: map[string]*schema.Schema{
			"variables": {
				Description: "The list of variables returned by the search",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabInstanceVariableGetSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabInstanceVariablesRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	options := &gitlab.ListInstanceVariablesOptions{
		Page:    1,
		PerPage: 20,
	}

	var variables []*gitlab.InstanceVariable
	for options.Page != 0 {
		paginatedVariables, resp, err := client.InstanceVariables.ListVariables(options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		variables = append(variables, paginatedVariables...)
		options.Page = resp.NextPage
	}

	d.SetId("instance_variables")
	if err := d.Set("variables", flattenGitlabInstanceVariables(variables)); err != nil {
		return diag.Errorf("failed to set variables to state: %v", err)
	}
	return nil
}

func flattenGitlabInstanceVariables(variables []*gitlab.InstanceVariable) (values []map[string]interface{}) {
	for _, variable := range variables {
		values = append(values, gitlabInstanceVariableToStateMap(variable))
	}
	return values
}
