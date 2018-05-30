package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataGitlabProject_basic(t *testing.T) {
	projectname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectConfig(projectname),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProject("gitlab_project.test", "data.gitlab_project.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProject(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		project := s.RootModule().Resources[src]
		projectResource := project.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchResource := search.Primary.Attributes

		if searchResource["id"] == "" {
			return fmt.Errorf("Expected to get a project ID from Gitlab")
		}

		testAttributes := []string{"id", "Name", "Path", "Visibility", "Description"}

		for _, attribute := range testAttributes {
			if searchResource[attribute] != projectResource[attribute] {
				return fmt.Errorf("Expected the project %s to be: %s, but got: %s", attribute, projectResource[attribute], searchResource[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabProjectConfig(projectname string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%s"
	path = "%s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
}

data "gitlab_project" "foo" {
	id = "${gitlab_project.test.id}"
}
	`, projectname, projectname)
}
