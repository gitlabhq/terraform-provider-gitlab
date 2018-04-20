package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabUsersRead,

		Schema: map[string]*schema.Schema{
			"options": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"order_by": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "id",
							ValidateFunc: validation.StringInSlice([]string{"id", "name",
								"username", "created_at"}, true),
						},
						"sort": {
							Type:         schema.TypeString,
							Optional:     true,
							Default:      "desc",
							ValidateFunc: validation.StringInSlice([]string{"desc", "asc"}, true),
						},
						"search": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"active": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"blocked": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"extern_uid": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"provider": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_before": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"created_after": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"username": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"email": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_admin": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"can_create_group": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"can_create_project": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"projects_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"extern_uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"organization": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"two_factor_enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceGitlabUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	listUsersOptions, id, err := expandGitlabUsersOptions(d.Get("options").([]interface{}))
	if err != nil {
		return err
	}
	log.Printf("\n\n\nListOptions\n%v\n\n\n", listUsersOptions)
	users, _, err := client.Users.ListUsers(listUsersOptions)
	log.Printf("\n\n\nListUsers\n%v\n\n\n", users)

	if err != nil {
		return err
	}

	d.Set("users", flattenGitlabUsers(users))
	d.SetId(fmt.Sprintf("%d", id))

	return nil
}

func flattenGitlabUsers(users []*gitlab.User) []interface{} {
	usersList := []interface{}{}

	for _, user := range users {
		values := map[string]interface{}{
			"user_id":            user.ID,
			"username":           user.Username,
			"email":              user.Email,
			"name":               user.Name,
			"is_admin":           user.IsAdmin,
			"can_create_group":   user.CanCreateGroup,
			"can_create_project": user.CanCreateProject,
			"projects_limit":     user.ProjectsLimit,
			"state":              user.State,
			"external":           user.External,
			"extern_uid":         user.ExternUID,
			"organization":       user.Organization,
			"two_factor_enabled": user.TwoFactorEnabled,
		}

		if user.CreatedAt != nil {
			values["created_at"] = user.CreatedAt
		}

		usersList = append(usersList, values)
	}

	return usersList
}

func expandGitlabUsersOptions(d []interface{}) (*gitlab.ListUsersOptions, int, error) {
	if len(d) == 0 {
		return nil, 0, nil
	}

	data := d[0].(map[string]interface{})
	listUsersOptions := &gitlab.ListUsersOptions{}
	options := ""

	if orderBy := data["order_by"].(string); orderBy != "" {
		listUsersOptions.OrderBy = &orderBy
		options += orderBy
	}
	if sort := data["sort"].(string); sort != "" {
		listUsersOptions.Sort = &sort
		options += sort
	}
	if search := data["search"].(string); search != "" {
		listUsersOptions.Search = &search
		options += search
	}
	if active := data["active"].(bool); active != false {
		listUsersOptions.Active = &active
		options += strconv.FormatBool(active)
	}
	if blocked := data["blocked"].(bool); blocked != false {
		listUsersOptions.Blocked = &blocked
		options += strconv.FormatBool(blocked)
	}
	if externalUID := data["extern_uid"].(string); externalUID != "" {
		listUsersOptions.ExternalUID = &externalUID
		options += externalUID
	}
	if provider := data["provider"].(string); provider != "" {
		// listUsersOptions.Provider = &provider
		options += provider
	}
	if createdBefore := data["created_before"].(string); createdBefore != "" {
		date, err := time.Parse("2006-01-02", createdBefore)
		if err != nil {
			return nil, 0, fmt.Errorf("created_before must be in yyyy-mm-dd format")
		}
		listUsersOptions.CreatedBefore = &date
		options += createdBefore
	}
	if createdAfter := data["created_after"].(string); createdAfter != "" {
		// date, err := time.Parse("2006-01-02", createdAfter)
		// if err != nil {
		// 	return nil, 0, fmt.Errorf("created_after must be in yyyy-mm-dd format")
		// }
		// listUsersOptions.CreatedAfter = &date
		options += createdAfter
	}

	id := schema.HashString(options)

	return listUsersOptions, id, nil
}
