package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_hooks", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_hooks`" + ` data source allows to retrieve details about hooks in a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#list-group-hooks)`,

		ReadContext: dataSourceGitlabGroupHooksRead,
		Schema: map[string]*schema.Schema{
			"group": {
				Description: "The ID or full path of the group.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hooks": {
				Description: "The list of hooks.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabGroupHookSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabGroupHooksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	group := d.Get("group").(string)
	options := gitlab.ListGroupHooksOptions{
		PerPage: 20,
		Page:    1,
	}

	var hooks []*gitlab.GroupHook
	for options.Page != 0 {
		paginatedHooks, resp, err := client.Groups.ListGroupHooks(group, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		hooks = append(hooks, paginatedHooks...)
		options.Page = resp.NextPage
	}

	d.SetId(group)
	if err := d.Set("hooks", flattenGitlabGroupHooks(group, hooks)); err != nil {
		return diag.Errorf("failed to set hooks to state: %v", err)
	}

	return nil
}

func flattenGitlabGroupHooks(group string, hooks []*gitlab.GroupHook) (values []map[string]interface{}) {
	for _, hook := range hooks {
		values = append(values, gitlabGroupHookToStateMap(group, hook))
	}
	return values
}
