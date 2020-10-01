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

	_ = d.Set("group_id", group.ID)
	_ = d.Set("full_path", group.FullPath)
	_ = d.Set("name", group.Name)
	_ = d.Set("full_name", group.FullName)
	_ = d.Set("web_url", group.WebURL)
	_ = d.Set("path", group.Path)
	_ = d.Set("description", group.Description)
	_ = d.Set("lfs_enabled", group.LFSEnabled)
	_ = d.Set("request_access_enabled", group.RequestAccessEnabled)
	_ = d.Set("visibility_level", group.Visibility)
	_ = d.Set("parent_id", group.ParentID)
	_ = d.Set("runners_token", group.RunnersToken)

	d.SetId(fmt.Sprintf("%d", group.ID))

	return nil
}
