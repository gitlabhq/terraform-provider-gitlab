package gitlab

import (
	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUsers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabUsersRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "all_users",
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
					},
				},
			},
		},
	}
}

func dataSourceGitlabUsersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	users, _, err := client.Users.ListUsers(nil)
	if err != nil {
		return err
	}

	d.Set("users", flattenGitlabUsers(users))
	d.SetId(d.Get("name").(string))

	return nil
}

func flattenGitlabUsers(users []*gitlab.User) []interface{} {
	usersList := []interface{}{}

	for _, user := range users {
		values := map[string]interface{}{
			"username":           user.Username,
			"email":              user.Email,
			"name":               user.Name,
			"is_admin":           user.IsAdmin,
			"can_create_group":   user.CanCreateGroup,
			"can_create_project": user.CanCreateProject,
			"projects_limit":     user.ProjectsLimit,
			"state":              user.State,
		}

		if user.CreatedAt != nil {
			values["created_at"] = user.CreatedAt.String()
		}

		usersList = append(usersList, values)
	}

	return usersList
}
