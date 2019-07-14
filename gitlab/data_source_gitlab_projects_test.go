package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataGitlabProjects_basic(t *testing.T) {
	projectname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectsConfig(projectname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_projects.foo", "projects.0.name", projectname),
				),
			},
		},
	})
}

func testAccDataGitlabProjectsConfig(projectname string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%s"
	path = "%s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
}

data "gitlab_projects" "foo" {
	search = "${gitlab_project.test.name}"
}
	`, projectname, projectname)
}
