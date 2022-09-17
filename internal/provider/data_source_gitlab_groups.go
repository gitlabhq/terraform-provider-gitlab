package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
	"strings"
)

var _ = registerDataSource("gitlab_groups", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_groups`" + ` data source allows details of multiple groups to be retrieved given some optional filter criteria.

-> Some attributes might not be returned depending on if you're an admin or not.

-> Some available options require administrator privileges.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/groups.html#list-groups)`,

		ReadContext: dataSourceGitlabGroupsRead,

		Schema: map[string]*schema.Schema{
			"order_by": {
				Description: "Order the groups' list by `id`, `name`, `path`, or `similarity`. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "name",
				ValidateFunc: validation.StringInSlice([]string{"id", "name",
					"path", "similarity"}, true),
			},
			"sort": {
				Description:  "Sort groups' list in asc or desc order. (Requires administrator privileges)",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "desc",
				ValidateFunc: validation.StringInSlice([]string{"desc", "asc"}, true),
			},
			"search": {
				Description: "Search groups by name or path.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"groups": {
				Description: "The list of groups.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Description: "The ID of the group.",
							Type:        schema.TypeInt,
							Computed:    true,
							Optional:    true,
						},
						"full_path": {
							Description: "The full path of the group.",
							Type:        schema.TypeString,
							Computed:    true,
							Optional:    true,
						},
						"name": {
							Description: "The name of this group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"full_name": {
							Description: "The full name of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"web_url": {
							Description: "Web URL of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"path": {
							Description: "The path of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"description": {
							Description: "The description of the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"lfs_enabled": {
							Description: "Boolean, is LFS enabled for projects in this group.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"request_access_enabled": {
							Description: "Boolean, is request for access enabled to the group.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
						"visibility_level": {
							Description: "Visibility level of the group. Possible values are `private`, `internal`, `public`.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"parent_id": {
							Description: "Integer, ID of the parent group.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"runners_token": {
							Description: "The group level registration token to use during runner setup.",
							Type:        schema.TypeString,
							Computed:    true,
							Sensitive:   true,
						},
						"default_branch_protection": {
							Description: "Whether developers and maintainers can push to the applicable default branch.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"prevent_forking_outside_group": {
							Description: "When enabled, users can not fork projects from this group to external namespaces.",
							Type:        schema.TypeBool,
							Computed:    true,
						},
					},
				},
			},
		},
	}
})

func dataSourceGitlabGroupsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	listGroupsOptions, id, err := expandGitlabGroupsOptions(d)
	if err != nil {
		return diag.FromErr(err)
	}
	page := 1
	groupslen := 0
	var groups []*gitlab.Group
	for page == 1 || groupslen != 0 {
		listGroupsOptions.Page = page
		paginatedGroups, _, err := client.Groups.ListGroups(listGroupsOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		groups = append(groups, paginatedGroups...)
		groupslen = len(paginatedGroups)
		page = page + 1
	}

	err = d.Set("groups", flattenGitlabGroups(groups))
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}

func flattenGitlabGroups(groups []*gitlab.Group) []interface{} {
	groupsList := []interface{}{}

	for _, group := range groups {
		values := map[string]interface{}{
			"group_id":                      group.ID,
			"full_path":                     group.FullPath,
			"name":                          group.Name,
			"full_name":                     group.Name,
			"web_url":                       group.WebURL,
			"path":                          group.Path,
			"description":                   group.Description,
			"lfs_enabled":                   group.LFSEnabled,
			"request_access_enabled":        group.RequestAccessEnabled,
			"visibility_level":              group.Visibility,
			"parent_id":                     group.ParentID,
			"runners_token":                 group.RunnersToken,
			"default_branch_protection":     group.DefaultBranchProtection,
			"prevent_forking_outside_group": group.PreventForkingOutsideGroup,
		}

		groupsList = append(groupsList, values)
	}

	return groupsList
}

func expandGitlabGroupsOptions(d *schema.ResourceData) (*gitlab.ListGroupsOptions, int, error) {
	listGroupsOptions := &gitlab.ListGroupsOptions{}
	var optionsHash strings.Builder

	if data, ok := d.GetOk("order_by"); ok {
		orderBy := data.(string)
		listGroupsOptions.OrderBy = &orderBy
		optionsHash.WriteString(orderBy)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("sort"); ok {
		sort := data.(string)
		listGroupsOptions.Sort = &sort
		optionsHash.WriteString(sort)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("search"); ok {
		search := data.(string)
		listGroupsOptions.Search = &search
		optionsHash.WriteString(search)
	}

	id := schema.HashString(optionsHash.String())

	return listGroupsOptions, id, nil
}
