package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabGroup() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabGroupRead,
		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"full_path",
				},
			},
			"full_path": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"group_id",
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"full_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"web_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lfs_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"request_access_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"visibility_level": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"parent_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"runners_token": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func dataSourceGitlabGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var group *gitlab.Group
	var err error

	log.Printf("[INFO] Reading Gitlab group")

	groupIDData, groupIDOk := d.GetOk("group_id")
	fullPathData, fullPathOk := d.GetOk("full_path")

	if groupIDOk {
		// Get group by id
		group, _, err = client.Groups.GetGroup(groupIDData.(int))
		if err != nil {
			return err
		}
	} else if fullPathOk {
		// Get group by full path
		group, _, err = client.Groups.GetGroup(fullPathData.(string))
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("one and only one of group_id or full_path must be set")
	}

	if err := setResourceData(d, map[string]interface{}{
		"group_id":               group.ID,
		"full_path":              group.FullPath,
		"name":                   group.Name,
		"full_name":              group.FullName,
		"web_url":                group.WebURL,
		"path":                   group.Path,
		"description":            group.Description,
		"lfs_enabled":            group.LFSEnabled,
		"request_access_enabled": group.RequestAccessEnabled,
		"visibility_level":       group.Visibility,
		"parent_id":              group.ParentID,
		"runners_token":          group.RunnersToken,
	}); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", group.ID))

	return nil
}
