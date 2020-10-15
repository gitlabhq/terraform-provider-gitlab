package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccDataSourceGitlabUser_basic(t *testing.T) {
	rString := fmt.Sprintf("%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			// Get user using its email
			{
				Config: testAccDataGitlabUserConfigEmail(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabUser("gitlab_user.foo", "data.gitlab_user.foo"),
				),
			},
			// Get user using its ID
			{
				Config: testAccDataGitlabUserConfigUserID(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabUser("gitlab_user.foo2", "data.gitlab_user.foo2"),
				),
			},
			// Get user using its username
			{
				Config: testAccDataGitlabUserConfigUsername(rString),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabUser("gitlab_user.foo", "data.gitlab_user.foo"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabUser(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		user := s.RootModule().Resources[src]
		userResource := user.Primary.Attributes

		search := s.RootModule().Resources[n]
		searchResource := search.Primary.Attributes

		testAttributes := []string{
			"username",
			"email",
			"name",
			"is_admin",
			"can_create_group",
			"projects_limit",
		}

		for _, attribute := range testAttributes {
			if searchResource[attribute] != userResource[attribute] {
				return fmt.Errorf("Expected user's parameter `%s` to be: %s, but got: `%s`", attribute, userResource[attribute], searchResource[attribute])
			}
		}

		return nil
	}
}

func testAccDataGitlabUserConfigEmail(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name     = "foo%s"
  username = "listest%s"
  password = "test%stt"
  email    = "listest%s@ssss.com"
  is_admin = false
}

resource "gitlab_user" "foo2" {
  name     = "foo2%s"
  username = "listest2%s"
  password = "test2%stt"
  email    = "listest2%s@ssss.com"
}

data "gitlab_user" "foo" {
  email = "${gitlab_user.foo.email}"
}
`, rString, rString, rString, rString, rString, rString, rString, rString)
}

func testAccDataGitlabUserConfigUserID(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name     = "foo%s"
  username = "listest%s"
  password = "test%stt"
  email    = "listest%s@ssss.com"
  is_admin = false
}

resource "gitlab_user" "foo2" {
  name     = "foo2%s"
  username = "listest2%s"
  password = "test2%stt"
  email    = "listest2%s@ssss.com"
}

data "gitlab_user" "foo2" {
  user_id = "${gitlab_user.foo2.id}"
}
`, rString, rString, rString, rString, rString, rString, rString, rString)
}

func testAccDataGitlabUserConfigUsername(rString string) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name     = "foo%s"
  username = "listest%s"
  password = "test%stt"
  email    = "listest%s@ssss.com"
  is_admin = false
}

resource "gitlab_user" "foo2" {
  name     = "foo2%s"
  username = "listest2%s"
  password = "test2%stt"
  email    = "listest2%s@ssss.com"
}

data "gitlab_user" "foo" {
  username = "${gitlab_user.foo.username}"
}
`, rString, rString, rString, rString, rString, rString, rString, rString)
}
