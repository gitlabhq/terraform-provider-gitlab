//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabProjectAccessToken_basic(t *testing.T) {
	project := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectAccessTokenDestroy,
		Steps: []resource.TestStep{
			// Create a basic access token.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_access_token" "foo" {
					project = %d
					name    = "foo"
					scopes  = ["api"]
				}
				`, project.ID),
				// Check computed and default attributes.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project_access_token.foo", "active", "true"),
					resource.TestCheckResourceAttr("gitlab_project_access_token.foo", "revoked", "false"),
					resource.TestCheckResourceAttr("gitlab_project_access_token.foo", "access_level", "maintainer"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "token"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "created_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "user_id"),
					resource.TestCheckNoResourceAttr("gitlab_project_access_token.foo", "expires_at"),
				),
			},
			// Verify upstream resource with an import.
			{
				ResourceName:      "gitlab_project_access_token.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// The token is only known during creating. We explicitly mention this limitation in the docs.
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Recreate the access token with updated attributes.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_access_token" "foo" {
					project = %d
					name    = "foo"
					scopes  = ["api", "read_api", "read_repository", "write_repository"]
					access_level = "developer"
					expires_at = %q
				}
				`, project.ID, time.Now().Add(time.Hour*48).Format("2006-01-02")),
				// Check computed and default attributes.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project_access_token.foo", "active", "true"),
					resource.TestCheckResourceAttr("gitlab_project_access_token.foo", "revoked", "false"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "token"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "created_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_access_token.foo", "user_id"),
				),
			},
			// Verify upstream resource with an import.
			{
				ResourceName:      "gitlab_project_access_token.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// The token is only known during creating. We explicitly mention this limitation in the docs.
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Recreate with `owner` access level.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_access_token" "foo" {
					project = %d
					name    = "foo"
					scopes  = ["api", "read_api", "read_repository", "write_repository"]
					access_level = "owner"
					expires_at = %q
				}
				`, project.ID, time.Now().Add(time.Hour*48).Format("2006-01-02")),
			},
			// Verify upstream resource with an import.
			{
				ResourceName:      "gitlab_project_access_token.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// The token is only known during creating. We explicitly mention this limitation in the docs.
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccCheckGitlabProjectAccessTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_access_token" {
			continue
		}

		project := rs.Primary.Attributes["project"]
		name := rs.Primary.Attributes["name"]

		tokens, _, err := testGitlabClient.ProjectAccessTokens.ListProjectAccessTokens(project, nil)
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if token.Name == name {
				return fmt.Errorf("project %q access token with name %q still exists", project, name)
			}
		}
	}

	return nil
}
