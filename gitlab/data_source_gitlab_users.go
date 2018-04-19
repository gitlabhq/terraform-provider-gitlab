package gitlab

import (
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
					},
				},
			},
		},
	}
}

func dataSourceGitlabUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	orderBy := d.Get("order_by").(string)
	sort := d.Get("sort").(string)
	listUsersOptions := &gitlab.ListUsersOptions{
		OrderBy: &orderBy,
		Sort:    &sort,
	}

	users, _, err := client.Users.ListUsers(listUsersOptions)
	if err != nil {
		return err
	}

	d.Set("users", flattenGitlabUsers(users))
	id := "all_users_" + orderBy + "_" + sort
	d.SetId(id)

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
		}

		if user.CreatedAt != nil {
			values["created_at"] = user.CreatedAt.String()
		}

		usersList = append(usersList, values)
	}

	return usersList
}
