package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func dataSourceGitlabProjectMembership() *schema.Resource {
	acceptedAccessLevels := make([]string, 0, len(accessLevelID))
	for k := range accessLevelID {
		acceptedAccessLevels = append(acceptedAccessLevels, k)
	}
	return &schema.Resource{
		Read: dataSourceGitlabProjectMembershipRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"inherited": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"access_level": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ValidateFunc: validateValueFunc(acceptedAccessLevels),
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

func dataSourceGitlabProjectMembershipRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var pm []*gitlab.ProjectMember
	var project *gitlab.Project
	var err error

	log.Printf("[INFO] Reading Gitlab group")

	projectIDData, _ := d.GetOk("id")
	inherited, _ := d.GetOk("inherited")

	// Get project by ID
	project, _, err = client.Projects.GetProject(projectIDData, nil)
	if err != nil {
		return err
	}

	log.Printf("[INFO] Reading Gitlab project memberships")

	// Get group memberships
	if inherited.(bool) {
		pm, _, err = client.ProjectMembers.ListAllProjectMembers(project.ID, &gitlab.ListProjectMembersOptions{})
		if err != nil {
			return err
		}
	} else {
		pm, _, err = client.ProjectMembers.ListProjectMembers(project.ID, &gitlab.ListProjectMembersOptions{})
		if err != nil {
			return err
		}
	}

	d.Set("members", flattenGitlabProjectMembers(d, pm))

	var optionsHash strings.Builder
	optionsHash.WriteString(strconv.Itoa(project.ID))

	if data, ok := d.GetOk("access_level"); ok {
		optionsHash.WriteString(data.(string))
	}

	id := schema.HashString(optionsHash.String())
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}

func flattenGitlabProjectMembers(d *schema.ResourceData, members []*gitlab.ProjectMember) []interface{} {
	membersList := []interface{}{}

	var filterAccessLevel gitlab.AccessLevelValue = gitlab.NoPermissions
	if data, ok := d.GetOk("access_level"); ok {
		filterAccessLevel = accessLevelID[data.(string)]
	}

	for _, member := range members {
		if filterAccessLevel != gitlab.NoPermissions && filterAccessLevel != member.AccessLevel {
			continue
		}

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
