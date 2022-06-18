//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabUsers_basic(t *testing.T) {
	rInt := acctest.RandInt()
	rInt2 := acctest.RandInt()
	user2 := fmt.Sprintf("user%d@test.test", rInt2)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabUsersConfig(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_user.foo", "name", "footest1"),
					resource.TestCheckResourceAttr("gitlab_user.foo2", "name", "footest2"),
				),
			},
			{
				Config: testAccDataSourceGitlabUsersConfigSort(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.email", user2),
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.projects_limit", "2"),
				),
			},
			{
				Config: testAccDataSourceGitlabUsersConfigSearch(rInt, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "1"),
					// resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.0.email", user2),
				),
			},
			{
				Config: testAccDataSourceGitlabLotsOfUsers(),
			},
			{
				Config: testAccDataSourceGitlabLotsOfUsersSearch(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_users.foo", "users.#", "99"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabUsersConfig(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabUsersConfigSort(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}

data "gitlab_users" "foo" {
  sort = "desc"
  search = "footest"
  order_by = "name"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabUsersConfigSearch(rInt int, rInt2 int) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = "footest1"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 3
}

resource "gitlab_user" "foo2" {
  name             = "footest2"
  username         = "listest%d"
  password         = "test%dtt"
  email            = "user%d@test.test"
  projects_limit   = 2
}

data "gitlab_users" "foo" {
  search = "user%d@test.test"
}
	`, rInt, rInt, rInt, rInt2, rInt2, rInt2, rInt2)
}

func testAccDataSourceGitlabLotsOfUsers() string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
  name             = format("lots user%%02d", count.index+1)
  username         = format("user%%02d", count.index+1)
  email            = format("user%%02d@example.com", count.index+1)
  password         = "8characters"
  count            = 99
}
`)
}

func testAccDataSourceGitlabLotsOfUsersSearch() string {
	return fmt.Sprintf(`%v
data "gitlab_users" "foo" {
	search = "lots"
}
	`, testAccDataSourceGitlabLotsOfUsers())
}
