package gitlab

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

// Search by email required
func dataSourceGitlabUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabUserRead,
		Schema: map[string]*schema.Schema{
			//Search option
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

// Performs the lookup
func dataSourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)
	// Create the query, grab the email for the query and set it for use
	var query *gitlab.ListUsersOptions
	email := strings.ToLower(d.Get("email").(string))
	*query.Search = email
	// Query to find the email. Returns a list
	users, _, err := client.Users.ListUsers(query)
	if err != nil {
		return err
	}

	// Create a user to save userdata to
	var user *gitlab.User
	// Grab User data out of list
	for _, a := range users {
		if a.Email == email {
			user = a
			break
		}
	}
	if user == nil {
		return fmt.Errorf("The email '%s' does not match any user email", email)
	}
	d.SetId(strconv.Itoa(user.ID))
	d.Set("name", user.Name)
	d.Set("username", user.Username)
	d.Set("email", user.Email)
	return nil
}
