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
					resource.TestCheckResourceAttr("gitlab_group.foo1", "name", fmt.Sprintf("foo1-name-%d", rInt)),
					resource.TestCheckResourceAttr("gitlab_group.foo2", "name", fmt.Sprintf("foo2-name-%d", rInt2)),
				),
			},
			{
				Config: testAccDataSourceGitlabGroupsConfigSort(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.0.name", fmt.Sprintf("foo1-name-%d", rInt)),
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.1.description", fmt.Sprintf("description-%d", rInt2)),
				),
			},
			{
				Config: testAccDataSourceGitlabGroupsConfigSearch(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.#", "1"),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroups(),
			},
			{
				Config: testAccDataSourceGitlabLotsOfGroupsSearch(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_groups.foos", "groups.#", "42"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabGroupsConfig(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo1" {
  name = "foo1-name-%d"
  path = "foo1-path-%d"
  description = "description-%d"
}

resource "gitlab_group" "foo2" {
  name = "foo2-name-%d"
  path = "foo2-path-%d"
  description = "description-%d"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabGroupsConfigSort(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo1" {
  name = "foo1-name-%d"
  path = "foo1-path-%d"
  description = "description-%d"
}

resource "gitlab_group" "foo2" {
  name = "foo2-name-%d"
  path = "foo2-path-%d"
  description = "description-%d"
}

data "gitlab_groups" "foos" {
  sort = "asc"
  search = "foo"
  order_by = "name"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabGroupsConfigSearch(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "description-%d"
}

resource "gitlab_group" "foo2" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "description-%d"
}

data "gitlab_groups" "foos" {
  search = "%d"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabLotsOfGroups() string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name             = format("lots-group-%%02d", count.index+1)
  description      = format("description-%%02d", count.index+1)
  path             = format("lots-group-path-%%02d", count.index+1)
  count            = 42
}
`)
}

func testAccDataSourceGitlabLotsOfGroupsSearch() string {
	return fmt.Sprintf(`%v
data "gitlab_groups" "foos" {
	search = "lots"
}
	`, testAccDataSourceGitlabLotsOfGroups())
}
