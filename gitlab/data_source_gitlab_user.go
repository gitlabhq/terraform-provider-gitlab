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
				Required: true,
			},
		},
	}
}

func dataSourceGitlabUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*gitlab.Client)

	log.Printf("[INFO] Reading Gitlab user")

	searchEmail := strings.ToLower(d.Get("email").(string))
	query := &gitlab.ListUsersOptions{
		Search: &searchEmail,
	}
	users, _, err := client.Users.ListUsers(query)
	if err != nil {
		return err
	}

	var found *gitlab.User

	for _, user := range users {
		if strings.ToLower(user.Email) == searchEmail {
			found = user
			break
		}
	}
	if found == nil {
		return fmt.Errorf("The email '%s' does not match any user email", searchEmail)
	}
	d.SetId(fmt.Sprintf("%d", found.ID))
	d.Set("name", found.Name)
	d.Set("username", found.Username)
	d.Set("email", found.Email)
	return nil
}
