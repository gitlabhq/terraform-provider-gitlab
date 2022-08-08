package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_hooks", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_hooks`" + ` data source allows to retrieve details about hooks in a project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/projects.html#list-project-hooks)`,

		ReadContext: dataSourceGitlabProjectHooksRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The name or id of the project.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"hooks": {
				Description: "The list of hooks.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabProjectHookSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabProjectHooksRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.ListProjectHooksOptions{
		PerPage: 20,
		Page:    1,
	}

	var hooks []*gitlab.ProjectHook
	for options.Page != 0 {
		paginatedHooks, resp, err := client.Projects.ListProjectHooks(project, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		hooks = append(hooks, paginatedHooks...)
		options.Page = resp.NextPage
	}

	d.SetId(project)
	if err := d.Set("hooks", flattenGitlabProjectHooks(project, hooks)); err != nil {
		return diag.Errorf("failed to set hooks to state: %v", err)
	}

	return nil
}

func flattenGitlabProjectHooks(project string, hooks []*gitlab.ProjectHook) (values []map[string]interface{}) {
	for _, hook := range hooks {
		values = append(values, gitlabProjectHookToStateMap(project, hook))
	}
	return values
}
