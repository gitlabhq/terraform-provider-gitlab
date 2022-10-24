package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_subgroups", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_subgroups`" + ` data source allows to get subgroups of a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#list-a-groups-subgroups)`,

		ReadContext: dataSourceGitlabGroupSubgroupsRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Required:    true,
			},
			"skip_groups": {
				Description: "Skip the group IDs passed.",
				Type:        schema.TypeList,
				Computed:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeInt},
			},
			"all_available": {
				Description: "Show all the groups you have access to.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"search": {
				Description: "Return the list of authorized groups matching the search criteria.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"order_by": {
				Description: "Order groups by name, path or id.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"sort": {
				Description: "Order groups in asc or desc order.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"statistics": {
				Description: "Include group statistics (administrators only).",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"with_custom_attributes": {
				Description: "Include custom attributes in response (administrators only).",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"owned": {
				Description: "Limit to groups explicitly owned by the current user.",
				Type:        schema.TypeBool,
				Computed:    true,
				Optional:    true,
			},
			"min_access_level": {
				Description: "Limit to groups where current user has at least this access level.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
			},
			"subgroups": {
				Description: "Subgroups of the parent group.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabGroupSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabGroupSubgroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var subgroups []*gitlab.Group
	var err error

	log.Printf("[INFO] Reading Gitlab group subgroups")

	groupIDData, groupIDOk := d.GetOk("group_id")

	if groupIDOk {
		// Get group subgroups by id
		subgroups, _, err = client.Groups.ListSubGroups(groupIDData.(int), nil, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("group_id is not valid")
	}

	d.SetId(fmt.Sprintf("%d", groupIDData))
	if err := d.Set("subgroups", flattenSubgroupsForState(subgroups)); err != nil {
		return diag.Errorf("Failed to set subgroups to state: %v", err)
	}

	return nil
}

func flattenSubgroupsForState(subgroups []*gitlab.Group) (values []map[string]interface{}) {
	for _, group := range subgroups {
		values = append(values, gitlabGroupToStateMap(group))
	}
	return values
}
