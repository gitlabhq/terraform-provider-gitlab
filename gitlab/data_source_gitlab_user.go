package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabUserRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"username",
					"email",
				},
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"user_id",
					"email",
				},
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"user_id",
					"username",
				},
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
			"note": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_provider": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bio": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"skype": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"linkedin": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"twitter": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"website_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"theme_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"color_scheme_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_sign_in_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"current_sign_in_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	var user *gitlab.User
	var err error

	log.Printf("[INFO] Reading Gitlab user")

	userIDData, userIDOk := d.GetOk("user_id")
	usernameData, usernameOk := d.GetOk("username")
	emailData, emailOk := d.GetOk("email")

	if userIDOk {
		// Get user by id
		user, _, err = client.Users.GetUser(userIDData.(int), gitlab.GetUsersOptions{})
		if err != nil {
			return err
		}
	} else if usernameOk || emailOk {
		username := strings.ToLower(usernameData.(string))
		email := strings.ToLower(emailData.(string))

		listUsersOptions := &gitlab.ListUsersOptions{}
		if usernameOk {
			// Get user by username
			listUsersOptions.Username = gitlab.String(username)
		} else {
			// Get user by email
			listUsersOptions.Search = gitlab.String(email)
		}

		var users []*gitlab.User
		users, _, err = client.Users.ListUsers(listUsersOptions)
		if err != nil {
			return err
		}

		if len(users) == 0 {
			return fmt.Errorf("couldn't find a user matching: %s%s", username, email)
		} else if len(users) != 1 {
			return fmt.Errorf("more than one user found matching: %s%s", username, email)
		}

		user = users[0]
	} else {
		return fmt.Errorf("one and only one of user_id, username or email must be set")
	}

	if err := setResourceData(d, map[string]interface{}{
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
		"created_at":         user.CreatedAt.String(),
		"organization":       user.Organization,
		"two_factor_enabled": user.TwoFactorEnabled,
		"note":               user.Note,
		"user_provider":      user.Provider,
		"avatar_url":         user.AvatarURL,
		"bio":                user.Bio,
		"location":           user.Location,
		"skype":              user.Skype,
		"linkedin":           user.Linkedin,
		"twitter":            user.Twitter,
		"website_url":        user.WebsiteURL,
		"theme_id":           user.ThemeID,
		"color_scheme_id":    user.ColorSchemeID,
	}); err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", user.ID))

	return nil
}
