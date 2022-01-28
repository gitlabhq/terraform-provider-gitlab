package gitlab

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUser() *schema.Resource {
	return &schema.Resource{
		Description: "Provide details about a specific user in the gitlab provider. Especially the ability to lookup the id for linking to other resources.\n\n" +
			"-> Some attributes might not be returned depending on if you're an admin or not. Please refer to [Gitlab documentation](https://docs.gitlab.com/ce/api/users.html#single-user) for more details.",

		ReadContext: dataSourceGitlabUserRead,
		Schema: map[string]*schema.Schema{
			"user_id": {
				Description: "The ID of the user.",
				Type:        schema.TypeInt,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"username",
					"email",
				},
			},
			"username": {
				Description: "The username of the user.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"user_id",
					"email",
				},
			},
			"email": {
				Description: "The email address of the user.",
				Type:        schema.TypeString,
				Computed:    true,
				Optional:    true,
				ConflictsWith: []string{
					"user_id",
					"username",
				},
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_admin": {
				Description: "Whether the user is an admin.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"can_create_group": {
				Description: "Whether the user can create groups.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"can_create_project": {
				Description: "Whether the user can create projects.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"projects_limit": {
				Description: "Number of projects the user can create.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"created_at": {
				Description: "Date the user was created at.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"state": {
				Description: "Whether the user is active or blocked.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"external": {
				Description: "Whether the user is external.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"extern_uid": {
				Description: "The external UID of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organization": {
				Description: "The organization of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"two_factor_enabled": {
				Description: "Whether user's two-factor auth is enabled.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"note": {
				Description: "Admin notes for this user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"user_provider": {
				Description: "The UID provider of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"avatar_url": {
				Description: "The avatar URL of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"bio": {
				Description: "The bio of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"location": {
				Description: "The location of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"skype": {
				Description: "Skype username of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"linkedin": {
				Description: "LinkedIn profile of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"twitter": {
				Description: "Twitter username of the user.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"website_url": {
				Description: "User's website URL.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"theme_id": {
				Description: "User's theme ID.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"color_scheme_id": {
				Description: "User's color scheme ID.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"last_sign_in_at": {
				Description: "Last user's sign-in date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"current_sign_in_at": {
				Description: "Current user's sign-in date.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func dataSourceGitlabUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	var user *gitlab.User
	var err error

	log.Printf("[INFO] Reading Gitlab user")

	userIDData, userIDOk := d.GetOk("user_id")
	usernameData, usernameOk := d.GetOk("username")
	emailData, emailOk := d.GetOk("email")

	if userIDOk {
		// Get user by id
		user, _, err = client.Users.GetUser(userIDData.(int), gitlab.GetUsersOptions{}, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
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
		users, _, err = client.Users.ListUsers(listUsersOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		if len(users) == 0 {
			return diag.Errorf("couldn't find a user matching: %s%s", username, email)
		} else if len(users) != 1 {
			return diag.Errorf("more than one user found matching: %s%s", username, email)
		}

		user = users[0]
	} else {
		return diag.Errorf("one and only one of user_id, username or email must be set")
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
	d.Set("created_at", user.CreatedAt.String())
	d.Set("organization", user.Organization)
	d.Set("two_factor_enabled", user.TwoFactorEnabled)
	d.Set("note", user.Note)
	d.Set("user_provider", user.Provider)
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
