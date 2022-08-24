//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabGroupHooks_basic(t *testing.T) {
	testAccCheckEE(t)

	testGroup := testAccCreateGroups(t, 1)[0]
	testHooks := testAccCreateGroupHooks(t, testGroup.ID, 25)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_group_hooks" "this" {
						group = "%s"
					}
				`, testGroup.FullPath),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_group_hooks.this", "hooks.#", fmt.Sprintf("%d", len(testHooks))),
					resource.TestCheckResourceAttr("data.gitlab_group_hooks.this", "hooks.0.url", testHooks[0].URL),
					resource.TestCheckResourceAttr("data.gitlab_group_hooks.this", "hooks.1.url", testHooks[1].URL),
				),
			},
		},
	})
}
