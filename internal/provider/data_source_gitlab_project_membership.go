package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_project_membership", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_project_membership`" + ` data source allows to list and filter all members of a project specified by either its id or full path.

-> **Note** exactly one of project_id or full_path must be provided.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/members.html#list-all-members-of-a-group-or-project)`,
		ReadContext: dataSourceGitlabProjectMembershipRead,
		Schema: map[string]*schema.Schema{
			"project_id": {
				Description:  "The ID of the project.",
				Type:         schema.TypeInt,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"project_id", "full_path"},
			},
			"full_path": {
				Description:  "The full path of the project.",
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ExactlyOneOf: []string{"project_id", "full_path"},
			},
			"query": {
				Description: "A query string to search for members",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"inherited": {
				Description: "Return all project members including members through ancestor groups",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"members": {
				Description: "The list of project members.",
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

func dataSourceGitlabProjectMembershipRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var project *gitlab.Project
	var err error

	log.Printf("[INFO] Reading Gitlab project")

	var pid interface{}
	if v, ok := d.GetOk("project_id"); ok {
		pid = v.(int)
	} else if v, ok := d.GetOk("full_path"); ok {
		pid = v.(string)
	} else {
		return diag.Errorf("one and only one of project_id or full_path must be set. This is a provider bug, please report upstream at https://github.com/gitlabhq/terraform-provider-gitlab/issues")
	}

	// Get project to have both, the `project_id` and `full_path` for setting the state.
	project, _, err = client.Projects.GetProject(pid, nil)
	if err != nil {
		return diag.FromErr(err)
	}

	var query *string
	if q, ok := d.GetOk("query"); ok {
		s := q.(string)
		query = &s
	}

	log.Printf("[INFO] Reading Gitlab project memberships")

	// Get project memberships
	listOptions := &gitlab.ListProjectMembersOptions{
		Query: query,
		ListOptions: gitlab.ListOptions{
			PerPage: 20,
			Page:    1,
		},
	}

	listMembers := client.ProjectMembers.ListProjectMembers
	if inherited, ok := d.GetOk("inherited"); ok && inherited.(bool) {
		listMembers = client.ProjectMembers.ListAllProjectMembers
	}

	var allPMs []*gitlab.ProjectMember
	for listOptions.Page != 0 {
		pms, resp, err := listMembers(project.ID, listOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		allPMs = append(allPMs, pms...)
		listOptions.Page = resp.NextPage
	}

	var optionsHash strings.Builder
	optionsHash.WriteString(strconv.Itoa(project.ID))

	if data, ok := d.GetOk("query"); ok {
		optionsHash.WriteString(data.(string))
	}

	id := schema.HashString(optionsHash.String())
	d.SetId(fmt.Sprintf("%d", id))

	d.Set("project_id", project.ID)
	d.Set("full_path", project.PathWithNamespace)

	if err := d.Set("members", flattenGitlabProjectMembers(d, allPMs)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func flattenGitlabProjectMembers(d *schema.ResourceData, members []*gitlab.ProjectMember) []interface{} {
	membersList := make([]interface{}, 0, len(members))
	for _, member := range members {
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
