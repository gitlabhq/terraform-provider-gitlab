package gitlab

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUsers() *schema.Resource {
	return &schema.Resource{
		Description: "Provide details about a list of users in the gitlab provider. The results include id, username, email, name and more about the requested users. Users can also be sorted and filtered using several options.\n\n" +
			"**NOTE**: Some available options require administrator privileges. Please visit [Gitlab API documentation][users_for_admins] for more information.",

		ReadContext: dataSourceGitlabUsersRead,

		Schema: map[string]*schema.Schema{
			"order_by": {
				Description: "Order the users' list by `id`, `name`, `username`, `created_at` or `updated_at`. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "id",
				ValidateFunc: validation.StringInSlice([]string{"id", "name",
					"username", "created_at", "updated_at"}, true),
			},
			"sort": {
				Description:  "Sort users' list in asc or desc order. (Requires administrator privileges)",
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "desc",
				ValidateFunc: validation.StringInSlice([]string{"desc", "asc"}, true),
			},
			"search": {
				Description: "Search users by username, name or email.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"active": {
				Description: "Filter users that are active.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"blocked": {
				Description: "Filter users that are blocked.",
				Type:        schema.TypeBool,
				Optional:    true,
			},
			"extern_uid": {
				Description: "Lookup users by external UID. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"extern_provider": {
				Description: "Lookup users by external provider. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created_before": {
				Description: "Search for users created before a specific date. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"created_after": {
				Description: "Search for users created after a specific date. (Requires administrator privileges)",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"users": {
				Description: "The list of users.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Description: "The unique id assigned to the user by the gitlab server.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"username": {
							Description: "The username of the user.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"email": {
							Description: "The e-mail address of the user.",
							Type:        schema.TypeString,
							Computed:    true,
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
						"provider": {
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
				},
			},
		},
	}
}

func dataSourceGitlabUsersRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	listUsersOptions, id, err := expandGitlabUsersOptions(d)
	if err != nil {
		return diag.FromErr(err)
	}
	page := 1
	userslen := 0
	var users []*gitlab.User
	for page == 1 || userslen != 0 {
		listUsersOptions.Page = page
		paginatedUsers, _, err := client.Users.ListUsers(listUsersOptions, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}
		users = append(users, paginatedUsers...)
		userslen = len(paginatedUsers)
		page = page + 1
	}

	d.Set("users", flattenGitlabUsers(users)) // lintignore: XR004 // TODO: Resolve this tfproviderlint issue
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
			"provider":           user.Provider,
			"two_factor_enabled": user.TwoFactorEnabled,
			"avatar_url":         user.AvatarURL,
			"bio":                user.Bio,
			"location":           user.Location,
			"skype":              user.Skype,
			"linkedin":           user.Linkedin,
			"twitter":            user.Twitter,
			"website_url":        user.WebsiteURL,
			"organization":       user.Organization,
			"theme_id":           user.ThemeID,
			"color_scheme_id":    user.ColorSchemeID,
		}

		if user.CreatedAt != nil {
			values["created_at"] = user.CreatedAt.String()
		}
		if user.LastSignInAt != nil {
			values["last_sign_in_at"] = user.LastSignInAt.String()
		}
		if user.CurrentSignInAt != nil {
			values["current_sign_in_at"] = user.CurrentSignInAt.String()
		}

		usersList = append(usersList, values)
	}

	return usersList
}

func expandGitlabUsersOptions(d *schema.ResourceData) (*gitlab.ListUsersOptions, int, error) {
	listUsersOptions := &gitlab.ListUsersOptions{}
	var optionsHash strings.Builder

	if data, ok := d.GetOk("order_by"); ok {
		orderBy := data.(string)
		listUsersOptions.OrderBy = &orderBy
		optionsHash.WriteString(orderBy)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("sort"); ok {
		sort := data.(string)
		listUsersOptions.Sort = &sort
		optionsHash.WriteString(sort)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("search"); ok {
		search := data.(string)
		listUsersOptions.Search = &search
		optionsHash.WriteString(search)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("active"); ok {
		active := data.(bool)
		listUsersOptions.Active = &active
		optionsHash.WriteString(strconv.FormatBool(active))
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("blocked"); ok {
		blocked := data.(bool)
		listUsersOptions.Blocked = &blocked
		optionsHash.WriteString(strconv.FormatBool(blocked))
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("extern_uid"); ok {
		externalUID := data.(string)
		listUsersOptions.ExternalUID = &externalUID
		optionsHash.WriteString(externalUID)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("extern_provider"); ok {
		provider := data.(string)
		listUsersOptions.Provider = &provider
		optionsHash.WriteString(provider)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("created_before"); ok {
		createdBefore := data.(string)
		date, err := time.Parse("2006-01-02", createdBefore)
		if err != nil {
			return nil, 0, fmt.Errorf("created_before must be in yyyy-mm-dd format")
		}
		listUsersOptions.CreatedBefore = &date
		optionsHash.WriteString(createdBefore)
	}
	optionsHash.WriteString(",")
	if data, ok := d.GetOk("created_after"); ok {
		createdAfter := data.(string)
		date, err := time.Parse("2006-01-02", createdAfter)
		if err != nil {
			return nil, 0, fmt.Errorf("created_after must be in yyyy-mm-dd format")
		}
		listUsersOptions.CreatedAfter = &date
		optionsHash.WriteString(createdAfter)
	}

	id := schema.HashString(optionsHash.String())

	return listUsersOptions, id, nil
}
