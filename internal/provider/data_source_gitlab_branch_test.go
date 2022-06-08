//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataGitlabBranch_basic(t *testing.T) {
	rInt := acctest.RandInt()
	project := testAccCreateProject(t)
	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabBranch(rInt, project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabBranch("gitlab_branch.foo", "data.gitlab_branch.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabBranch(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		branch := s.RootModule().Resources[src]
		branchAttr := branch.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchAttr := search.Primary.Attributes

		testAttributes := []string{
			"id",
			"name",
			"web_url",
			"default",
			"project",
			"can_push",
			"merged",
			"commit",
			"parent_ids",
			"protected",
			"developer_can_merge",
			"developer_can_push",
		}

		for _, attribute := range testAttributes {
			if searchAttr[attribute] != branchAttr[attribute] {
				return fmt.Errorf("expected branch's parameter `%s` to be: %s, but got: `%s`", attribute, branchAttr[attribute], searchAttr[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabBranch(rInt int, project string) string {
	return fmt.Sprintf(`
%s

data "gitlab_branch" "foo" {
  name = "${gitlab_branch.foo.name}"
  project = "%s"
}
`, testAccDataGitlabBranchSetup(rInt, project), project)
}

func testAccDataGitlabBranchSetup(rInt int, project string) string {
	return fmt.Sprintf(`
resource "gitlab_branch" "foo" {
	name = "testbranch-%[1]d"
	ref = "main"
	project = "%s"
}
  `, rInt, project)
}
