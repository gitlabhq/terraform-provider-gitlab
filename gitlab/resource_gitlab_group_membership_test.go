package gitlab

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabGroupMembership_basic(t *testing.T) {
	var membership gitlab.GroupMember
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupMembershipDestroy,
		Steps: []resource.TestStep{

			// Assign member to the group as a developer
			{
				Config: testAccGitlabGroupMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabGroupMembershipExists("gitlab_group_membership.foo", &membership), testAccCheckGitlabGroupMembershipAttributes(&membership, &testAccGitlabGroupMembershipExpectedAttributes{
					access_level: fmt.Sprintf("developer"),
				})),
			},

			// Update the group member to change the access level (use testAccGitlabGroupMembershipUpdateConfig for Config)
			{
				Config: testAccGitlabGroupMembershipUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabGroupMembershipExists("gitlab_group_membership.foo", &membership), testAccCheckGitlabGroupMembershipAttributes(&membership, &testAccGitlabGroupMembershipExpectedAttributes{
					access_level: fmt.Sprintf("guest"),
				})),
			},

			// Update the group member to change the access level back
			{
				Config: testAccGitlabGroupMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabGroupMembershipExists("gitlab_group_membership.foo", &membership), testAccCheckGitlabGroupMembershipAttributes(&membership, &testAccGitlabGroupMembershipExpectedAttributes{
					access_level: fmt.Sprintf("developer"),
				})),
			},
		},
	})
}

func testAccCheckGitlabGroupMembershipExists(n string, membership *gitlab.GroupMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		conn := testAccProvider.Meta().(*gitlab.Client)
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("No group ID is set")
		}

		userID := rs.Primary.Attributes["user_id"]
		id, _ := strconv.Atoi(userID)
		if userID == "" {
			return fmt.Errorf("No user id is set")
		}

		gotGroupMembership, _, err := conn.GroupMembers.GetGroupMember(groupID, id)
		if err != nil {
			return err
		}

		*membership = *gotGroupMembership
		return nil
	}
}

type testAccGitlabGroupMembershipExpectedAttributes struct {
	access_level string
}

func testAccCheckGitlabGroupMembershipAttributes(membership *gitlab.GroupMember, want *testAccGitlabGroupMembershipExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		access_level_id, ok := accessLevel[membership.AccessLevel]
		if !ok {
			return fmt.Errorf("Invalid access level '%s'", access_level_id)
		}
		if access_level_id != want.access_level {
			return fmt.Errorf("got access level %s; want %s", access_level_id, want.access_level)
		}
		return nil
	}
}

func testAccCheckGitlabGroupMembershipDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group_membership" {
			continue
		}

		groupID := rs.Primary.Attributes["group_id"]
		userID := rs.Primary.Attributes["user_id"]

		// GetGroupMember needs int type for userID
		userIDI, err := strconv.Atoi(userID)
		gotMembership, resp, err := conn.GroupMembers.GetGroupMember(groupID, userIDI)
		if err != nil {
			if gotMembership != nil && fmt.Sprintf("%d", gotMembership.AccessLevel) == rs.Primary.Attributes["access_level"] {
				return fmt.Errorf("Group still has member.")
			}
			return nil
		}

		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabGroupMembershipConfig(rInt int) string {
	return fmt.Sprintf(`resource "gitlab_group_membership" "foo" {
group_id = "${gitlab_group.foo.id}"
user_id = "${gitlab_user.test.id}"
access_level = "developer"
}

resource "gitlab_group" "foo" {
name = "foo%d"
path = "foo%d"
}

resource "gitlab_user" "test" {
name = "foo%d"
username = "listest%d"
password = "test%dtt"
email = "listest%d@ssss.com"
}
`, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabGroupMembershipUpdateConfig(rInt int) string {
	return fmt.Sprintf(`resource "gitlab_group_membership" "foo" {
group_id = "${gitlab_group.foo.id}"
user_id = "${gitlab_user.test.id}"
access_level = "guest"
}

resource "gitlab_group" "foo" {
name = "foo%d"
path = "foo%d"
}

resource "gitlab_user" "test" {
name = "foo%d"
username = "listest%d"
password = "test%dtt"
email = "listest%d@ssss.com"
}
`, rInt, rInt, rInt, rInt, rInt)
}
