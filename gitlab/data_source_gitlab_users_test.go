package gitlab

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccDataSourceGitlabUsers_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccComputeLbIpRangesConfig,
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr("data.gitlab_users.test",
						"users.#", regexp.MustCompile("^[1-9]*$"))),
			},
		},
	})
}

const testAccComputeLbIpRangesConfig = `
data "gitlab_users" "test" {}
`
