//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabPersonalAccessToken_basic(t *testing.T) {
	user := testAccCreateUsers(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabPersonalAccessTokenDestroy,
		Steps: []resource.TestStep{
			// Create a basic access token.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_personal_access_token" "foo" {
					user_id = %d
					name    = "foo"
					scopes  = ["api"]
				}
				`, user.ID),
				// Check computed and default attributes.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "active", "true"),
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "revoked", "false"),
					resource.TestCheckResourceAttrSet("gitlab_personal_access_token.foo", "token"),
					resource.TestCheckResourceAttrSet("gitlab_personal_access_token.foo", "created_at"),
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "user_id", fmt.Sprintf("%d", user.ID)),
					resource.TestCheckNoResourceAttr("gitlab_personal_access_token.foo", "expires_at"),
				),
			},
			// Verify upstream resource with an import.
			{
				ResourceName:      "gitlab_personal_access_token.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// The token is only known during creating. We explicitly mention this limitation in the docs.
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Recreate the access token with updated attributes.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_personal_access_token" "foo" {
					user_id    = %d
					name       = "foo"
					scopes     = ["api", "read_user", "read_api", "read_repository", "write_repository", "sudo", "read_registry", "write_registry"]
					expires_at = %q
				}
				`, user.ID, time.Now().Add(time.Hour*48).Format("2006-01-02")),
				// Check computed and default attributes.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "active", "true"),
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "revoked", "false"),
					resource.TestCheckResourceAttrSet("gitlab_personal_access_token.foo", "token"),
					resource.TestCheckResourceAttrSet("gitlab_personal_access_token.foo", "created_at"),
					resource.TestCheckResourceAttr("gitlab_personal_access_token.foo", "user_id", fmt.Sprintf("%d", user.ID)),
				),
			},
			// Verify upstream resource with an import.
			{
				ResourceName:      "gitlab_personal_access_token.foo",
				ImportState:       true,
				ImportStateVerify: true,
				// The token is only known during creating. We explicitly mention this limitation in the docs.
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccCheckGitlabPersonalAccessTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_personal_access_token" {
			continue
		}

		name := rs.Primary.Attributes["name"]
		userId := rs.Primary.Attributes["user_id"]

		userIdInt, err := strconv.Atoi(userId)
		if err != nil {
			return fmt.Errorf("Error converting user ID to string: %v", userId)
		}

		tokens, _, err := testGitlabClient.PersonalAccessTokens.ListPersonalAccessTokens(&gitlab.ListPersonalAccessTokensOptions{UserID: &userIdInt})
		if err != nil {
			return err
		}

		for _, token := range tokens {
			if token.Name == name && !token.Revoked {
				return fmt.Errorf("personal access token with name %q is not in a revoked state", name)
			}
		}
	}

	return nil
}
