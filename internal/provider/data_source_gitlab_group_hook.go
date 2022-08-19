package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_hook", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_hook`" + ` data source allows to retrieve details about a hook in a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#get-group-hook)`,

		ReadContext: dataSourceGitlabGroupHookRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabGroupHookSchema(), []string{"group", "hook_id"}, nil),
	}
})

func dataSourceGitlabGroupHookRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	group := d.Get("group").(string)
	hookID := d.Get("hook_id").(int)

	hook, _, err := client.Groups.GetGroupHook(group, hookID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", group, hookID))
	stateMap := gitlabGroupHookToStateMap(group, hook)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
