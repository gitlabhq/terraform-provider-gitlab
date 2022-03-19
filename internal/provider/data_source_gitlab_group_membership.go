package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_group_membership", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_group_membership`" + ` data source allows to list and filter all members of a group specified by either its id or full path.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/members.html#list-all-members-of-a-group-or-project)`,

		ReadContext: dataSourceGitlabGroupMembershipRead,
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
			"access_level": {
				Description:      "Only return members with the desired access level. Acceptable values are: `guest`, `reporter`, `developer`, `maintainer`, `owner`.",
				Type:             schema.TypeString,
				Computed:         true,
				Optional:         true,
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validGroupAccessLevelNames, false)),
			},
			"members": {
				Description: "The list of group members.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The unique id assigned to the user by the gitlab server.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"username": {
							Description: "The username of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"name": {
							Description: "The name of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"state": {
							Description: "Whether the user is active or blocked.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"avatar_url": {
							Description: "The avatar URL of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"web_url": {
							Description: "User's website URL.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"access_level": {
							Description: "The level of access to the group.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"expires_at": {
							Description: "Expiration date for the group membership.",
							Type:        schema.TypeString,
							Computed:    true,
						},
					},
				},
			},
		},
	}
})

func dataSourceGitlabGroupMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	log.Printf("[INFO] Reading Gitlab group memberships")

	// Get group memberships
	listOptions := &gitlab.ListGroupMembersOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	var allGms []*gitlab.GroupMember
	for {
		gms, resp, err := client.Groups.ListGroupMembers(group.ID, listOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		allGms = append(allGms, gms...)

		if resp.NextPage == 0 {
			break
		}
		listOptions.Page = resp.NextPage
	}

	d.Set("group_id", group.ID)
	d.Set("full_path", group.FullPath)

	d.Set("members", flattenGitlabMembers(d, allGms)) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue

	var optionsHash strings.Builder
	optionsHash.WriteString(strconv.Itoa(group.ID))

	if data, ok := d.GetOk("access_level"); ok {
		optionsHash.WriteString(data.(string))
	}

	id := schema.HashString(optionsHash.String())
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}

func flattenGitlabMembers(d *schema.ResourceData, members []*gitlab.GroupMember) []interface{} {
	membersList := []interface{}{}

	var filterAccessLevel gitlab.AccessLevelValue = gitlab.NoPermissions
	if data, ok := d.GetOk("access_level"); ok {
		filterAccessLevel = accessLevelNameToValue[data.(string)]
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
			"access_level": accessLevelValueToName[gitlab.AccessLevelValue(member.AccessLevel)],
		}

		if member.ExpiresAt != nil {
			values["expires_at"] = member.ExpiresAt.String()
		}

		membersList = append(membersList, values)
	}

	return membersList
}
