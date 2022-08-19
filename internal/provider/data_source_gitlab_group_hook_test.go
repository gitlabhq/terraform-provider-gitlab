//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabGroupHook_basic(t *testing.T) {
	testAccCheckEE(t)

	testGroup := testAccCreateGroups(t, 1)[0]
	testHook := testAccCreateGroupHooks(t, testGroup.ID, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_group_hook" "this" {
						group   = "%s"
						hook_id = %d
					}
				`, testGroup.FullPath, testHook.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_group_hook.this", "hook_id", fmt.Sprintf("%d", testHook.ID)),
					resource.TestCheckResourceAttr("data.gitlab_group_hook.this", "group_id", fmt.Sprintf("%d", testGroup.ID)),
					resource.TestCheckResourceAttr("data.gitlab_group_hook.this", "url", testHook.URL),
				),
			},
		},
	})
}
