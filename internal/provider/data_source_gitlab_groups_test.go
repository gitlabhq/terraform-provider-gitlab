//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabGroups_basic(t *testing.T) {
	prefixFoo := "acctest-group-foo"
	groupsFoo := testAccCreateGroupsWithPrefix(t, 2, prefixFoo)

	prefixLotsOf := "acctest-group-lotsof"
	testAccCreateGroupsWithPrefix(t, 42, prefixLotsOf)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabGroupsConfigSearchSort(prefixFoo),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.name", groupsFoo[0].Name),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroupsSearch(prefixLotsOf),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.lotsof", "groups.#", "42"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabGroupsConfigSearchSort(prefix string) string {
	return fmt.Sprintf(`
data "gitlab_groups" "foos" {
  sort = "asc"
  search = "%s"
  order_by = "id"
}
	`, prefix)
}

func testAccDataSourceGitlabLotsOfGroupsSearch(prefix string) string {
	return fmt.Sprintf(`
data "gitlab_groups" "lotsof" {
	search = "%s"
}
	`, prefix)
}
