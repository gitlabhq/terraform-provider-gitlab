package gitlab

import (
	"fmt"
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
			"identities_extern_uid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"identities_provider": {
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
			"users": {
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

	listUsersOptions, id, err := expandGitlabUsersOptions(d)
	if err != nil {
		return err
	}
	users, _, err := client.Users.ListUsers(listUsersOptions)

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
			"id":                 user.ID,
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
			values["created_at"] = user.CreatedAt.String()
		}

		usersList = append(usersList, values)
	}

	return usersList
}

func expandGitlabUsersOptions(d *schema.ResourceData) (*gitlab.ListUsersOptions, int, error) {
	listUsersOptions := &gitlab.ListUsersOptions{}
	optionsHash := ""

	if data, ok := d.GetOk("order_by"); ok {
		orderBy := data.(string)
		listUsersOptions.OrderBy = &orderBy
		optionsHash += orderBy
	}
	optionsHash += ","
	if data, ok := d.GetOk("sort"); ok {
		sort := data.(string)
		listUsersOptions.Sort = &sort
		optionsHash += sort
	}
	optionsHash += ","
	if data, ok := d.GetOk("search"); ok {
		search := data.(string)
		listUsersOptions.Search = &search
		optionsHash += search
	}
	optionsHash += ","
	if data, ok := d.GetOk("active"); ok {
		active := data.(bool)
		listUsersOptions.Active = &active
		optionsHash += strconv.FormatBool(active)
	}
	optionsHash += ","
	if data, ok := d.GetOk("blocked"); ok {
		blocked := data.(bool)
		listUsersOptions.Blocked = &blocked
		optionsHash += strconv.FormatBool(blocked)
	}
	optionsHash += ","
	if data, ok := d.GetOk("identities_extern_uid"); ok {
		externalUID := data.(string)
		listUsersOptions.ExternalUID = &externalUID
		optionsHash += externalUID
	}
	optionsHash += ","
	if data, ok := d.GetOk("identities_provider"); ok {
		provider := data.(string)
		listUsersOptions.Provider = &provider
		optionsHash += provider
	}
	optionsHash += ","
	if data, ok := d.GetOk("created_before"); ok {
		createdBefore := data.(string)
		date, err := time.Parse("2006-01-02", createdBefore)
		if err != nil {
			return nil, 0, fmt.Errorf("created_before must be in yyyy-mm-dd format")
		}
		listUsersOptions.CreatedBefore = &date
		optionsHash += createdBefore
	}
	optionsHash += ","
	if data, ok := d.GetOk("created_after"); ok {
		createdAfter := data.(string)
		date, err := time.Parse("2006-01-02", createdAfter)
		if err != nil {
			return nil, 0, fmt.Errorf("created_after must be in yyyy-mm-dd format")
		}
		listUsersOptions.CreatedAfter = &date
		optionsHash += createdAfter
	}

	id := schema.HashString(optionsHash)

	return listUsersOptions, id, nil
}
