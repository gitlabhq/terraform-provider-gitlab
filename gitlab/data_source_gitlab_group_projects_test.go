package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataGitlabGroupProjects_basic(t *testing.T) {
	projectname := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabGroupProjectsConfig(projectname),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_group_projects.foo", "projects.0.name", projectname),
				),
			},
		},
	})
}

func testAccDataGitlabGroupProjectsConfig(projectname string) string {
	return fmt.Sprintf(`
resource "gitlab_project" "test"{
	name = "%s"
	path = "%s"
	description = "Terraform acceptance tests"
	visibility_level = "public"
	namespace_id = "${gitlab_group.foo.id}"
}

resource "gitlab_group" "foo" {
	name = "test-%s"
	path = "test-%s"

	visibility_level = "public"
}


data "gitlab_group_projects" "foo" {
	search = "${gitlab_project.test.name}"
	group_id = "${gitlab_group.foo.id}"
}
	`, projectname, projectname, projectname, projectname)
}
