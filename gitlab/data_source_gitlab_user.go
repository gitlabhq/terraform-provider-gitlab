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

	d.Set("user_id", user.ID)
	d.Set("username", user.Username)
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("is_admin", user.IsAdmin)
	d.Set("can_create_group", user.CanCreateGroup)
	d.Set("can_create_project", user.CanCreateProject)
	d.Set("projects_limit", user.ProjectsLimit)
	d.Set("state", user.State)
	d.Set("external", user.External)
	d.Set("extern_uid", user.ExternUID)
	d.Set("created_at", user.CreatedAt)
	d.Set("organization", user.Organization)
	d.Set("two_factor_enabled", user.TwoFactorEnabled)
	d.Set("note", user.Note)
	d.Set("provider", user.Provider)
	d.Set("avatar_url", user.AvatarURL)
	d.Set("bio", user.Bio)
	d.Set("location", user.Location)
	d.Set("skype", user.Skype)
	d.Set("linkedin", user.Linkedin)
	d.Set("twitter", user.Twitter)
	d.Set("website_url", user.WebsiteURL)
	d.Set("theme_id", user.ThemeID)
	d.Set("color_scheme_id", user.ColorSchemeID)

	d.SetId(fmt.Sprintf("%d", user.ID))

	return nil
}
