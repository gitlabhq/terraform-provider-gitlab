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

		srcR := s.RootModule().Resources[src]
		srcA := srcR.Primary.Attributes

		r := s.RootModule().Resources[n]
		a := r.Primary.Attributes

		if a["id"] == "" {
			return fmt.Errorf("Expected to get a user email from Gitlab")
		}

		testAtts := []string{"id", "email"}

		for _, att := range testAtts {
			if a[att] != srcA[att] {
				return fmt.Errorf("Expected the user %s to be: %s, but got: %s", att, srcA[att], a[att])
			}
		}
		return nil
	}
}

func testAccDataGitlabUserConfig(userEmail string) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo"{
	email = "%s"
}

data "gitlab_user" "foo" {
	name = "${gitlab_user.foo.email}"
}
	`, userEmail)
}
