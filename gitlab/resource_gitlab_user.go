package gitlab

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceGitlabUserCreate,
		Read:   resourceGitlabUserRead,
		Update: resourceGitlabUserUpdate,
		Delete: resourceGitlabUserDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{

			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
				ConflictsWith: []string{
					"user_id",
				},
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"email": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
				ConflictsWith: []string{
					"user_id",
				},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_create_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"can_create_project": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"projects_limit": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"extern_uid": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"two_factor_enabled": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"user_provider": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "provider.gitlab",
			},
			"avatar_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bio": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"location": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"skype": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"linkedin": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"twitter": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"website_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"theme_id": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
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
			"private_profile": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_confirmation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceGitlabUserSetToState(d *schema.ResourceData, user *gitlab.User) {
	d.Set("user_id", user.ID)
	d.Set("username", user.Username)
	d.Set("email", user.Email)
	d.Set("name", user.Name)
	d.Set("is_admin", user.IsAdmin)
	d.Set("is_external", user.External)
	d.Set("can_create_group", user.CanCreateGroup)
	d.Set("can_create_project", user.CanCreateProject)
	d.Set("projects_limit", user.ProjectsLimit)
	d.Set("created_at", user.CreatedAt)
	d.Set("state", user.State)
	d.Set("extern_uid", user.ExternUID)
	d.Set("organization", user.Organization)
	d.Set("two_factor_enabled", user.TwoFactorEnabled)
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
	d.Set("last_sign_in_at", user.LastSignInAt)
	d.Set("current_sign_in_at", user.CurrentSignInAt)
	d.Set("private_profile", user.PrivateProfile)
}

func resourceGitlabUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateUserOptions{
		Username:         gitlab.String(d.Get("username").(string)),
		Password:         gitlab.String(d.Get("password").(string)),
		Email:            gitlab.String(d.Get("email").(string)),
		Name:             gitlab.String(d.Get("name").(string)),
		Admin:            gitlab.Bool(d.Get("is_admin").(bool)),
		External:         gitlab.Bool(d.Get("is_external").(bool)),
		CanCreateGroup:   gitlab.Bool(d.Get("can_create_group").(bool)),
		ProjectsLimit:    gitlab.Int(d.Get("projects_limit").(int)),
		ExternUID:        gitlab.String(d.Get("extern_uid").(string)),
		Organization:     gitlab.String(d.Get("organization").(string)),
		Provider:         gitlab.String(d.Get("user_provider").(string)),
		Bio:              gitlab.String(d.Get("bio").(string)),
		Location:         gitlab.String(d.Get("location").(string)),
		Skype:            gitlab.String(d.Get("skype").(string)),
		Linkedin:         gitlab.String(d.Get("linkedin").(string)),
		Twitter:          gitlab.String(d.Get("twitter").(string)),
		WebsiteURL:       gitlab.String(d.Get("website_url").(string)),
		PrivateProfile:   gitlab.Bool(d.Get("private_profile").(bool)),
		SkipConfirmation: gitlab.Bool(d.Get("skip_confirmation").(bool)),
	}

	log.Printf("[DEBUG] create gitlab user %q", *options.Username)

	user, _, err := client.Users.CreateUser(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", user.ID))
	d.Set("is_admin", user.IsAdmin)
	d.Set("is_external", user.External)

	return resourceGitlabUserRead(d, meta)
}

func resourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] read gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	user, _, err := client.Users.GetUser(id)
	if err != nil {
		log.Printf("[DEBUG] Error for user %s", d.Id())
		return err
	}

	resourceGitlabUserSetToState(d, user)
	return nil
}

func resourceGitlabUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.ModifyUserOptions{}

	if d.HasChange("username") {
		options.Username = gitlab.String(d.Get("username").(string))
	}

	if d.HasChange("password") {
		options.Password = gitlab.String(d.Get("password").(string))
	}

	if d.HasChange("email") {
		options.Email = gitlab.String(d.Get("email").(string))
	}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("is_admin") {
		options.Admin = gitlab.Bool(d.Get("is_admin").(bool))
	}

	if d.HasChange("is_external") {
		options.Admin = gitlab.Bool(d.Get("is_external").(bool))
	}

	if d.HasChange("can_create_group") {
		options.CanCreateGroup = gitlab.Bool(d.Get("can_create_group").(bool))
	}

	if d.HasChange("projects_limit") {
		options.ProjectsLimit = gitlab.Int(d.Get("projects_limit").(int))
	}

	if d.HasChange("extern_uid") {
		options.ExternUID = gitlab.String(d.Get("extern_uid").(string))
	}

	if d.HasChange("organization") {
		options.Organization = gitlab.String(d.Get("organization").(string))
	}

	if d.HasChange("user_provider") {
		options.Provider = gitlab.String(d.Get("user_provider").(string))
	}

	if d.HasChange("bio") {
		options.Bio = gitlab.String(d.Get("bio").(string))
	}

	if d.HasChange("location") {
		options.Twitter = gitlab.String(d.Get("location").(string))
	}

	if d.HasChange("skype") {
		options.Skype = gitlab.String(d.Get("skype").(string))
	}

	if d.HasChange("linkedin") {
		options.Linkedin = gitlab.String(d.Get("linkedin").(string))
	}

	if d.HasChange("twitter") {
		options.Twitter = gitlab.String(d.Get("twitter").(string))
	}

	if d.HasChange("website_url") {
		options.WebsiteURL = gitlab.String(d.Get("website_url").(string))
	}

	if d.HasChange("private_profile") {
		options.PrivateProfile = gitlab.Bool(d.Get("private_profile").(bool))
	}

	log.Printf("[DEBUG] update gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	_, _, err := client.Users.ModifyUser(id, options)
	if err != nil {
		return err
	}

	return resourceGitlabUserRead(d, meta)
}

func resourceGitlabUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	_, err := client.Users.DeleteUser(id)
	// Ignoring error due to some bug in library
	log.Printf("[DEBUG] Delete gitlab user %s", err)
	return nil
}
