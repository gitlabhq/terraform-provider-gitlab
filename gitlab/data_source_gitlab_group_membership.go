package gitlab

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabGroupMembership() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabGroupMembershipRead,
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
			"members": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"avatar_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"web_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"expires_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabGroupMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var gm []*gitlab.GroupMember
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

	log.Printf("[INFO] Reading Gitlab group memberships")

	// Get group memberships
	gm, _, err = client.Groups.ListGroupMembers(group.ID, &gitlab.ListGroupMembersOptions{})
	if err != nil {
		return err
	}

	d.Set("group_id", group.ID)
	d.Set("full_path", group.FullPath)

	d.Set("members", flattenGitlabMembers(gm))
	d.SetId(fmt.Sprintf("%d", group.ID))

	return nil
}

func flattenGitlabMembers(members []*gitlab.GroupMember) []interface{} {
	membersList := []interface{}{}

	for _, member := range members {
		values := map[string]interface{}{
			"id":           member.ID,
			"username":     member.Username,
			"name":         member.Name,
			"state":        member.State,
			"avatar_url":   member.AvatarURL,
			"web_url":      member.WebURL,
			"access_level": accessLevel[gitlab.AccessLevelValue(member.AccessLevel)],
		}

		if member.ExpiresAt != nil {
			values["expires_at"] = member.ExpiresAt.String()
		}

		membersList = append(membersList, values)
	}

	return membersList
}
