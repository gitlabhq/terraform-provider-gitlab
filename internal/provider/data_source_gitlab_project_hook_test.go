//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabProjectHook_basic(t *testing.T) {
	testProject := testAccCreateProject(t)
	testHook := testAccCreateProjectHooks(t, testProject.ID, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_project_hook" "this" {
						project = "%s"
						hook_id = %d
					}
				`, testProject.PathWithNamespace, testHook.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_project_hook.this", "hook_id", fmt.Sprintf("%d", testHook.ID)),
					resource.TestCheckResourceAttr("data.gitlab_project_hook.this", "project_id", fmt.Sprintf("%d", testProject.ID)),
					resource.TestCheckResourceAttr("data.gitlab_project_hook.this", "url", testHook.URL),
				),
			},
		},
	})
}
