package gitlab

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_admin": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"can_create_group": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"skip_confirmation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"projects_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  0,
			},
			"is_external": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"reset_password": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"note": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceGitlabUserSetToState(d *schema.ResourceData, user *gitlab.User) {
	d.Set("username", user.Username)
	d.Set("name", user.Name)
	d.Set("can_create_group", user.CanCreateGroup)
	d.Set("projects_limit", user.ProjectsLimit)
	d.Set("email", user.Email)
	d.Set("is_admin", user.IsAdmin)
	d.Set("is_external", user.External)
	d.Set("note", user.Note)
	d.Set("skip_confirmation", user.ConfirmedAt != nil && !user.ConfirmedAt.IsZero())
}

func resourceGitlabUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	options := &gitlab.CreateUserOptions{
		Email:            gitlab.String(d.Get("email").(string)),
		Password:         gitlab.String(d.Get("password").(string)),
		Username:         gitlab.String(d.Get("username").(string)),
		Name:             gitlab.String(d.Get("name").(string)),
		ProjectsLimit:    gitlab.Int(d.Get("projects_limit").(int)),
		Admin:            gitlab.Bool(d.Get("is_admin").(bool)),
		CanCreateGroup:   gitlab.Bool(d.Get("can_create_group").(bool)),
		SkipConfirmation: gitlab.Bool(d.Get("skip_confirmation").(bool)),
		External:         gitlab.Bool(d.Get("is_external").(bool)),
		ResetPassword:    gitlab.Bool(d.Get("reset_password").(bool)),
		Note:             gitlab.String(d.Get("note").(string)),
	}

	if *options.Password == "" && !*options.ResetPassword {
		return fmt.Errorf("At least one of either password or reset_password must be defined")
	}

	log.Printf("[DEBUG] create gitlab user %q", *options.Username)

	user, _, err := client.Users.CreateUser(options)
	if err != nil {
		return err
	}

	d.SetId(fmt.Sprintf("%d", user.ID))

	return resourceGitlabUserRead(d, meta)
}

func resourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] import -- read gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	user, _, err := client.Users.GetUser(id, gitlab.GetUsersOptions{})
	if err != nil {
		return err
	}

	resourceGitlabUserSetToState(d, user)
	return nil
}

func resourceGitlabUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	options := &gitlab.ModifyUserOptions{}

	if d.HasChange("name") {
		options.Name = gitlab.String(d.Get("name").(string))
	}

	if d.HasChange("username") {
		options.Username = gitlab.String(d.Get("username").(string))
	}

	if d.HasChange("email") {
		options.Email = gitlab.String(d.Get("email").(string))
		options.SkipReconfirmation = gitlab.Bool(true)
	}

	if d.HasChange("is_admin") {
		options.Admin = gitlab.Bool(d.Get("is_admin").(bool))
	}

	if d.HasChange("can_create_group") {
		options.CanCreateGroup = gitlab.Bool(d.Get("can_create_group").(bool))
	}

	if d.HasChange("projects_limit") {
		options.ProjectsLimit = gitlab.Int(d.Get("projects_limit").(int))
	}

	if d.HasChange("is_external") {
		options.External = gitlab.Bool(d.Get("is_external").(bool))
	}

	if d.HasChange("note") {
		options.Note = gitlab.String(d.Get("note").(string))
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

	if _, err := client.Users.DeleteUser(id); err != nil {
		return err
	}

	stateConf := &resource.StateChangeConf{
		Timeout: 5 * time.Minute,
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			user, resp, err := client.Users.GetUser(id, gitlab.GetUsersOptions{})
			if resp != nil && resp.StatusCode == 404 {
				return user, "Deleted", nil
			}
			if err != nil {
				return user, "Error", err
			}
			return user, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("Could not finish deleting user %d: %w", id, err)
	}

	return nil
}
