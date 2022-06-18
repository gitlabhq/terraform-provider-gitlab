//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabMembership_basic(t *testing.T) {
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Create the group and one member
			{
				Config: testAccDataSourceGitlabGroupMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_group.foo", "name", fmt.Sprintf("foo%d", rInt)),
					resource.TestCheckResourceAttr("gitlab_user.test", "name", fmt.Sprintf("foo%d", rInt)),
					resource.TestCheckResourceAttr("gitlab_group_membership.foo", "access_level", "developer"),
				),
			},
			{
				Config: testAccDataSourceGitlabGroupMembershipConfig_basic(rInt),
				Check: resource.ComposeTestCheckFunc(
					// Members is 2 because the user owning the token is always added to the group
					resource.TestCheckResourceAttr("data.gitlab_group_membership.foo", "members.#", "2"),
					resource.TestCheckResourceAttr("data.gitlab_group_membership.foo", "members.1.username", fmt.Sprintf("listest%d", rInt)),
				),
			},

			// Get group using its ID, but return maintainers only
			{
				Config: testAccDataSourceGitlabGroupMembershipConfigFilterAccessLevel(rInt),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_group_membership.foomaintainers", "members.#", "0"),
				),
			},
		},
	})
}

func TestAccDataSourceGitlabMembership_pagination(t *testing.T) {
	userCount := 21

	group := testAccCreateGroups(t, 1)[0]
	users := testAccCreateUsers(t, userCount)
	testAccAddGroupMembers(t, group.ID, users)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataSourceGitlabGroupMembershipPagination(group.ID),
				Check:  resource.TestCheckResourceAttr("data.gitlab_group_membership.this", "members.#", fmt.Sprintf("%d", userCount)),
			},
		},
	})
}

func testAccDataSourceGitlabGroupMembershipConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo%d"
  path = "foo%d"
}

resource "gitlab_user" "test" {
  name     = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email    = "listest%d@ssss.com"
}

resource "gitlab_group_membership" "foo" {
  group_id     = "${gitlab_group.foo.id}"
  user_id      = "${gitlab_user.test.id}"
  access_level = "developer"
}`, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccDataSourceGitlabGroupMembershipConfig_basic(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo%d"
  path = "foo%d"
}

data "gitlab_group_membership" "foo" {
  group_id = "${gitlab_group.foo.id}"
}`, rInt, rInt)
}

func testAccDataSourceGitlabGroupMembershipConfigFilterAccessLevel(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo%d"
  path = "foo%d"
}

data "gitlab_group_membership" "foomaintainers" {
  group_id     = "${gitlab_group.foo.id}"
  access_level = "maintainer"
}`, rInt, rInt)
}

func testAccDataSourceGitlabGroupMembershipPagination(groupId int) string {
	return fmt.Sprintf(`
data "gitlab_group_membership" "this" {
  group_id     = "%d"
  access_level = "developer"
}`, groupId)
}
