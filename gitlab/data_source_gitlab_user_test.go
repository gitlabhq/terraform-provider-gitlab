package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataGitlabUser_basic(t *testing.T) {
	userEmail := fmt.Sprintf("tf-%s", acctest.RandString(5))

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabUserConfig(userEmail),
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

		if searchResource["email"] == "" {
			return fmt.Errorf("Expected to get user email from Gitlab")
		}

		testAttributes := []string{"email"}

		for _, attribute := range testAttributes {
			if searchResource[attribute] != userResource[attribute] {
				return fmt.Errorf("Expected the user %s to be: %s, but got: %s", attribute, userResource[attribute], searchResource[attribute])
			}
		}
		return nil
	}
}

func testAccDataGitlabUserConfig(userEmail string) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "foo %s"
  username         = "listest%s"
  password         = "test%stt"
  email            = "listest%s@ssss.com"
}

data "gitlab_user" "foo" {
	email = "${gitlab_user.foo.email}"
}
	`, userEmail, userEmail, userEmail, userEmail)
}
