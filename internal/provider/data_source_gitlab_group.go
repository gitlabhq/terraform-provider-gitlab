package provider

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Provide details about a specific group in the gitlab provider.\n\n" +
			"> **Note**: exactly one of group_id or full_path must be provided.",

		ReadContext: dataSourceGitlabGroupRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Description: "The ID of the group.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"full_path",
				},
			},
			"full_path": {
				Description: "The full path of the group.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"group_id",
				},
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
		},
	}
}

func dataSourceGitlabGroupRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var group *gitlab.Group
	var err error

	log.Printf("[INFO] Reading Gitlab group")

	groupIDData, groupIDOk := d.GetOk("group_id")
	fullPathData, fullPathOk := d.GetOk("full_path")

	if groupIDOk {
		// Get group by id
		group, _, err = client.Groups.GetGroup(groupIDData.(int), nil, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else if fullPathOk {
		// Get group by full path
		group, _, err = client.Groups.GetGroup(fullPathData.(string), nil, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		return diag.Errorf("one and only one of group_id or full_path must be set")
	}

	d.Set("group_id", group.ID)
	d.Set("full_path", group.FullPath)
	d.Set("name", group.Name)
	d.Set("full_name", group.FullName)
	d.Set("web_url", group.WebURL)
	d.Set("path", group.Path)
	d.Set("description", group.Description)
	d.Set("lfs_enabled", group.LFSEnabled)
	d.Set("request_access_enabled", group.RequestAccessEnabled)
	d.Set("visibility_level", group.Visibility)
	d.Set("parent_id", group.ParentID)
	d.Set("runners_token", group.RunnersToken)
	d.Set("default_branch_protection", group.DefaultBranchProtection)

	d.SetId(fmt.Sprintf("%d", group.ID))

	return nil
}
