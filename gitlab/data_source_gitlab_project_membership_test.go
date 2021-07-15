package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccDataSourceGitlabProjectMembership_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Create the project and one member
			{
				Config: testAccDataSourceGitlabProjectMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project.foo", "name", fmt.Sprintf("foo%d", rInt)),
					resource.TestCheckResourceAttr("gitlab_user.test", "name", fmt.Sprintf("foo%d", rInt)),
					resource.TestCheckResourceAttr("gitlab_project_membership.foo", "access_level", "developer"),
				),
			},
			{
				Config: testAccDataSourceGitlabProjectMembershipConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					// Members is 2 because the user owning the token is always added to the project
					resource.TestCheckResourceAttr("data.gitlab_project_membership.foo", "members.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_project_membership.foo", "members.1.username", fmt.Sprintf("listest%d", rInt)),
				),
			},

			// Get project using its ID, but return reporters only
			{
				Config: testAccDataSourceGitlabProjectMembershipConfigFilterAccessLevel(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_project_membership.fooreporters", "members.#", "0"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectMembershipConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo%d"
  path = "foo%d"
}

resource "gitlab_user" "test" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}

resource "gitlab_project_membership" "foo" {
  project_id   = "${gitlab_project.foo.id}"
  user_id      = "${gitlab_user.test.id}"
  access_level = "developer"
}`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccDataSourceGitlabProjectMembershipConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo%d"
  path = "foo%d"
}

data "gitlab_project_membership" "foo" {
  id = "${gitlab_project.foo.id}"
}`, rInt, rInt)
}

func testAccDataSourceGitlabProjectMembershipConfigFilterAccessLevel(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo%d"
  path = "foo%d"
}

data "gitlab_project_membership" "fooreporters" {
  id           = "${gitlab_project.foo.id}"
  access_level = "owner"
}`, rInt, rInt)
}
