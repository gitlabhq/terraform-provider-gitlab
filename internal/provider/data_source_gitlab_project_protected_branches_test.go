//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataGitlabProjectProtectedBranches_search(t *testing.T) {
	projectName := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectProtectedBranchesConfigGetProjectSearch(projectName),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(
						"data.gitlab_project_protected_branches.test",
						"protected_branches.0.name",
						"main",
					),
					resource.TestCheckResourceAttr(
						"data.gitlab_project_protected_branches.test",
						"protected_branches.0.push_access_levels.0.access_level",
						"maintainer",
					),
				),
			},
		},
	})
}

func testAccDataGitlabProjectProtectedBranchesConfigGetProjectSearch(projectName string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test" {
  name           = "%s"
  path           = "%s"
  default_branch = "main"
}

resource "gitlab_branch_protection" "test" {
  project            = gitlab_project.test.id
  branch             = "main"
  push_access_level  = "maintainer"
  merge_access_level = "developer"
}

data "gitlab_project_protected_branches" "test" {
  project_id = gitlab_branch_protection.test.project # This expresses the dependency of the data source on the protected branch having first been configured
}
`, projectName, projectName)
}
