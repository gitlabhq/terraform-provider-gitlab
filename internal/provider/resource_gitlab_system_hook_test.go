//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabSystemHook_basic(t *testing.T) {
	var hook gitlab.Hook
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabSystemHookDestroy,
		Steps: []resource.TestStep{
			// Create a hook with all options
			{
				Config: testAccGitlabSystemHookConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabSystemHookExists("gitlab_system_hook.this", &hook),
					resource.TestCheckResourceAttrSet("gitlab_system_hook.this", "created_at"),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_system_hook.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Update the hook to toggle all the values to their inverse
			{
				Config: testAccGitlabSystemHookUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabSystemHookExists("gitlab_system_hook.this", &hook),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_system_hook.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccCheckGitlabSystemHookExists(n string, hook *gitlab.Hook) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		hookID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotHook, _, err := testGitlabClient.SystemHooks.GetHook(hookID)
		if err != nil {
			return err
		}
		*hook = *gotHook
		return nil
	}
}

func testAccCheckGitlabSystemHookDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_system_hook" {
			continue
		}
		hookID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotHook, _, err := testGitlabClient.SystemHooks.GetHook(hookID, nil)
		if err == nil {
			if gotHook != nil && gotHook.ID == hookID {
				return fmt.Errorf("System Hook %d still exists after deletion", hookID)
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabSystemHookConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_system_hook" "this" {
  url                      = "https://example.com/hook-%d"
  token                    = "secret-token"
  push_events              = true
  tag_push_events          = true
  merge_requests_events    = true
  repository_update_events = true
  enable_ssl_verification  = true
}
	`, rInt)
}

func testAccGitlabSystemHookUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_system_hook" "this" {
  url                      = "https://example.com/hook-%d"
  token                    = "another-secret-token"
  push_events              = false
  tag_push_events          = false
  merge_requests_events    = false
  repository_update_events = false
  enable_ssl_verification  = false
}
	`, rInt)
}
