package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

var validUserStateValues = []string{
	"active",
	"deactivated",
	"blocked",
}

var _ = registerResource("gitlab_user", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_user`" + ` resource allows to manage the lifecycle of a user.

-> the provider needs to be configured with admin-level access for this resource to work.

-> You must specify either password or reset_password.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/users.html)`,

		CreateContext: resourceGitlabUserCreate,
		ReadContext:   resourceGitlabUserRead,
		UpdateContext: resourceGitlabUserUpdate,
		DeleteContext: resourceGitlabUserDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"username": {
				Description: "The username of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"password": {
				Description: "The password of the user.",
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				ForceNew:    true,
			},
			"email": {
				Description: "The e-mail address of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"name": {
				Description: "The name of the user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"is_admin": {
				Description: "Boolean, defaults to false.  Whether to enable administrative privileges",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"can_create_group": {
				Description: "Boolean, defaults to false. Whether to allow the user to create groups.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"skip_confirmation": {
				Description: "Boolean, defaults to true. Whether to skip confirmation.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
			},
			"projects_limit": {
				Description: "Integer, defaults to 0.  Number of projects user can create.",
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     0,
			},
			"is_external": {
				Description: "Boolean, defaults to false. Whether a user has access only to some internal or private projects. External users can only access projects to which they are explicitly granted access.",
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
			},
			"reset_password": {
				Description: "Boolean, defaults to false. Send user password reset link.",
				Type:        schema.TypeBool,
				Optional:    true,
				ForceNew:    true,
			},
			"note": {
				Description: "The note associated to the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"state": {
				Description:      fmt.Sprintf("String, defaults to 'active'. The state of the user account. Valid values are %s.", renderValueListForDocs(validUserStateValues)),
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "active",
				ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validUserStateValues, false)),
			},
			"namespace_id": {
				Description: "The ID of the user's namespace. Available since GitLab 14.10.",
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
			},
		},
	}
})

func resourceGitlabUserSetToState(d *schema.ResourceData, user *gitlab.User) {
	d.Set("username", user.Username)
	d.Set("name", user.Name)
	d.Set("can_create_group", user.CanCreateGroup)
	d.Set("projects_limit", user.ProjectsLimit)
	d.Set("email", user.Email)
	d.Set("is_admin", user.IsAdmin)
	d.Set("is_external", user.External)
	d.Set("note", user.Note)
	d.Set("state", user.State)
	d.Set("namespace_id", user.NamespaceID)
}

func resourceGitlabUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return diag.Errorf("At least one of either password or reset_password must be defined")
	}

	log.Printf("[DEBUG] create gitlab user %q", *options.Username)

	user, _, err := client.Users.CreateUser(options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d", user.ID))

	if d.Get("state") == "blocked" {
		err := client.Users.BlockUser(user.ID, gitlab.WithContext(ctx))

		if err != nil {
			return diag.FromErr(err)
		}
	} else if d.Get("state") == "deactivated" {
		err := client.Users.DeactivateUser(user.ID, gitlab.WithContext(ctx))

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabUserRead(ctx, d, meta)
}

func resourceGitlabUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] import -- read gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	user, _, err := client.Users.GetUser(id, gitlab.GetUsersOptions{}, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] gitlab user not found %d", id)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	resourceGitlabUserSetToState(d, user)
	return nil
}

func resourceGitlabUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	_, _, err := client.Users.ModifyUser(id, options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("state") {
		oldState, newState := d.GetChange("state")
		var err error
		// NOTE: yes, this can be written much more consice, however, for the sake of understanding the behavior,
		//       of the API and the allowed state transitions of GitLab, let's keep it as-is and enjoy the readability.
		if newState == "active" && oldState == "blocked" {
			err = client.Users.UnblockUser(id, gitlab.WithContext(ctx))
		} else if newState == "active" && oldState == "deactivated" {
			err = client.Users.ActivateUser(id, gitlab.WithContext(ctx))
		} else if newState == "blocked" && oldState == "active" {
			err = client.Users.BlockUser(id, gitlab.WithContext(ctx))
		} else if newState == "blocked" && oldState == "deactivated" {
			err = client.Users.BlockUser(id, gitlab.WithContext(ctx))
		} else if newState == "deactivated" && oldState == "active" {
			err = client.Users.DeactivateUser(id, gitlab.WithContext(ctx))
		} else if newState == "deactivated" && oldState == "blocked" {
			// a blocked user cannot be deactivated, GitLab will return an error, like:
			// `403 Forbidden - A blocked user cannot be deactivated by the API`
			// we have to unblock the user first
			err = client.Users.UnblockUser(id, gitlab.WithContext(ctx))
			if err != nil {
				return diag.FromErr(err)
			}
			err = client.Users.DeactivateUser(id, gitlab.WithContext(ctx))
		}

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceGitlabUserRead(ctx, d, meta)
}

func resourceGitlabUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	log.Printf("[DEBUG] Delete gitlab user %s", d.Id())

	id, _ := strconv.Atoi(d.Id())

	if _, err := client.Users.DeleteUser(id, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	stateConf := &resource.StateChangeConf{
		Timeout: 5 * time.Minute,
		Target:  []string{"Deleted"},
		Refresh: func() (interface{}, string, error) {
			user, resp, err := client.Users.GetUser(id, gitlab.GetUsersOptions{}, gitlab.WithContext(ctx))
			if resp != nil && resp.StatusCode == 404 {
				return user, "Deleted", nil
			}
			if err != nil {
				return user, "Error", err
			}
			return user, "Deleting", nil
		},
	}

	if _, err := stateConf.WaitForStateContext(ctx); err != nil {
		return diag.Errorf("Could not finish deleting user %d: %s", id, err)
	}

	return nil
}
