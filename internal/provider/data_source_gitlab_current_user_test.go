//go:build acceptance
// +build acceptance

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/xanzy/go-gitlab"
)

func TestAccDataSourceGitlabCurrentUser_basic(t *testing.T) {
	//The root user has no public email by default, set the public email so it shows up properly.
	testGitlabClient.Users.ModifyUser(1, &gitlab.ModifyUserOptions{
		Email: gitlab.String("admin@example.com"),
		// The public email MUST match an email on record for the user, or it gets a bad request.
		PublicEmail: gitlab.String("admin@example.com"),
	})

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "gitlab_current_user" "this" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "id"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "name"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "username"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "bot"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "group_count"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "namespace_id"),
					resource.TestCheckResourceAttrSet("data.gitlab_current_user.this", "public_email"),
				),
			},
		},
	})
}
