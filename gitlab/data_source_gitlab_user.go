package gitlab

import (
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func dataSourceGitlabUser() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGitlabUserRead,
		Schema: map[string]*schema.Schema{
			"email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dataSourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab user")

	searchEmail := strings.ToLower(d.Get("email").(string))
	userName := strings.ToLower(d.Get("username").(string))
	var q *string
	if searchEmail != "" {
		q = &searchEmail
	} else {
		q = &userName
	}
	query := &gitlab.ListUsersOptions{
		Search: q,
	}
	users, _, err := client.Users.ListUsers(query)
	if err != nil {
		return err
	}

	var found *gitlab.User

	if searchEmail != "" {
		for _, user := range users {
			if strings.ToLower(user.Email) == searchEmail {
				found = user
				break
			}
		}
	} else {
		for _, user := range users {
			if strings.ToLower(user.Username) == userName {
				found = user
				break
			}
		}
	}

	if found == nil {
		return fmt.Errorf("The email '%s' does not match any user email", searchEmail)
	}
	d.SetId(fmt.Sprintf("%d", found.ID))
	d.Set("name", found.Name)
	d.Set("userName", found.Username)
	d.Set("email", found.Email)
	return nil
}
