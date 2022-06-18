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

func TestAccDataGitlabProjectTag_basic(t *testing.T) {
	rInt := acctest.RandInt()
	project := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectTag(rInt, project.PathWithNamespace),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjectTag("gitlab_project_tag.foo", "data.gitlab_project_tag.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectTag(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		tag := s.RootModule().Resources[src]
		tagAttr := tag.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchAttr := search.Primary.Attributes

		testAttributes := []string{
			"id",
			"name",
			"project",
			"message",
			"protected",
			"target",
			"release",
			"commit",
		}

		for _, attribute := range testAttributes {
			if searchAttr[attribute] != tagAttr[attribute] {
				return fmt.Errorf("expected the parameter of tag `%s` to be: %s, but got: `%s`", attribute, tagAttr[attribute], searchAttr[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabProjectTag(rInt int, project string) string {
	return fmt.Sprintf(`
%s
data "gitlab_project_tag" "foo" {
  name    = "${gitlab_project_tag.foo.name}"
  project = "%s"
}
`, testAccDataGitlabProjectTagSetup(rInt, project), project)
}

func testAccDataGitlabProjectTagSetup(rInt int, project string) string {
	return fmt.Sprintf(`
resource "gitlab_project_tag" "foo" {
    name    = "tag-%[1]d"
    ref     = "main"
    project = "%s"
}
  `, rInt, project)
}
