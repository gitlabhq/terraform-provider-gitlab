//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabGroups_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rInt2 := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabGroupsConfig(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_group.foo", "name", "footest1"),
					resource.TestCheckResourceAttr("gitlab_group.foo2", "name", "footest2"),
				),
			},
			{
				Config: testAccDataSourceGitlabGroupsConfigSort(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foo", "groups.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_groups.foo", "groups.0.name", "footest1"),
					resource.TestCheckResourceAttr("data.gitlab_groups.foo", "groups.1.description", "description2"),
				),
			},
			{
				Config: testAccDataSourceGitlabGroupsConfigSearch(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foo", "groups.#", "1"),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroups(),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroupsSearch(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foo", "groups.#", "99"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabGroupsConfig(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}

resource "gitlab_group" "foo2" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}
	`, rInt, rInt, rInt2, rInt2)
}

func testAccDataSourceGitlabGroupsConfigSort(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}

resource "gitlab_group" "foo2" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}

data "gitlab_groups" "foo" {
  sort = "desc"
  search = "footest"
  order_by = "name"
}
	`, rInt, rInt, rInt2, rInt2)
}

func testAccDataSourceGitlabGroupsConfigSearch(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}

resource "gitlab_group" "foo2" {
  name        = "footest%d"
  description = "description%d"
  path        = "/"
}

data "gitlab_groups" "foo" {
  search = "footest%d"
}
	`, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabLotsOfGroups() string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name             = format("lots group%%02d", count.index+1)
  description      = format("description%%02d", count.index+1)
  path             = "/"
  count            = 99
}
`)
}

func testAccDataSourceGitlabLotsOfGroupsSearch() string {
	return fmt.Sprintf(`%v
data "gitlab_groups" "foo" {
	search = "lots"
}
	`, testAccDataSourceGitlabLotsOfGroups())
}
