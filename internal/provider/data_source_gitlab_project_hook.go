package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_hook`" + ` data source allows to retrieve details about a hook in a project.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/projects.html#get-project-hook)`,

		ReadContext: dataSourceGitlabProjectHookRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabProjectHookSchema(), []string{"project", "hook_id"}, nil),
	}
})

func dataSourceGitlabProjectHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project := d.Get("project").(string)
	hookID := d.Get("hook_id").(int)

	hook, _, err := client.Projects.GetProjectHook(project, hookID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", project, hookID))
	stateMap := gitlabProjectHookToStateMap(project, hook)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
