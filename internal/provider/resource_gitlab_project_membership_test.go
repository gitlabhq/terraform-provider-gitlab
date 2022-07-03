//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectMembership_basic(t *testing.T) {
	var membership gitlab.ProjectMember
	rInt := acctest.RandInt()

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectMembershipDestroy,
		Steps: []resource.TestStep{

			// Assign member to the project as a developer
			{
				Config: testAccGitlabProjectMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabProjectMembershipExists("gitlab_project_membership.foo", &membership), testAccCheckGitlabProjectMembershipAttributes(&membership, &testAccGitlabProjectMembershipExpectedAttributes{
					access_level: "developer",
				})),
			},

			// Update the project member to change the access level (use testAccGitlabProjectMembershipUpdateConfig for Config)
			{
				Config: testAccGitlabProjectMembershipUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabProjectMembershipExists("gitlab_project_membership.foo", &membership), testAccCheckGitlabProjectMembershipAttributes(&membership, &testAccGitlabProjectMembershipExpectedAttributes{
					access_level: "guest",
					expiresAt:    "2099-01-01",
				})),
			},

			// Update the project member to change the access level back
			{
				Config: testAccGitlabProjectMembershipConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabProjectMembershipExists("gitlab_project_membership.foo", &membership), testAccCheckGitlabProjectMembershipAttributes(&membership, &testAccGitlabProjectMembershipExpectedAttributes{
					access_level: "developer",
				})),
			},
		},
	})
}

func testAccCheckGitlabProjectMembershipExists(n string, membership *gitlab.ProjectMember) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("No project ID is set")
		}

		userID := rs.Primary.Attributes["user_id"]
		id, _ := strconv.Atoi(userID)
		if userID == "" {
			return fmt.Errorf("No user id is set")
		}

		gotProjectMembership, _, err := testGitlabClient.ProjectMembers.GetProjectMember(projectID, id)
		if err != nil {
			return err
		}

		*membership = *gotProjectMembership
		return nil
	}
}

type testAccGitlabProjectMembershipExpectedAttributes struct {
	access_level string
	expiresAt    string
}

func testAccCheckGitlabProjectMembershipAttributes(membership *gitlab.ProjectMember, want *testAccGitlabProjectMembershipExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		access_level_id, ok := accessLevelValueToName[membership.AccessLevel]
		if !ok {
			return fmt.Errorf("Invalid access level '%s'", access_level_id)
		}
		if access_level_id != want.access_level {
			return fmt.Errorf("got access level %s; want %s", access_level_id, want.access_level)
		}
		return nil
	}
}

func testAccCheckGitlabProjectMembershipDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_membership" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		userID := rs.Primary.Attributes["user_id"]

		// GetProjectMember needs int type for userID
		userIDI, err := strconv.Atoi(userID) // nolint // TODO: Resolve this golangci-lint issue: ineffectual assignment to err (ineffassign)
		gotMembership, _, err := testGitlabClient.ProjectMembers.GetProjectMember(projectID, userIDI)
		if err != nil {
			if gotMembership != nil && fmt.Sprintf("%d", gotMembership.AccessLevel) == rs.Primary.Attributes["access_level"] {
				return fmt.Errorf("Project still has member.")
			}
			return nil
		}

		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectMembershipConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_membership" "foo" {
  project_id = "${gitlab_project.foo.id}"
  user_id = "${gitlab_user.test.id}"
  access_level = "developer"
}

resource "gitlab_project" "foo" {
  name = "foo%d"
  description = "Terraform acceptance tests"
  visibility_level ="public"
}

resource "gitlab_user" "test" {
  name = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email = "listest%d@ssss.com"
}
`, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabProjectMembershipUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_membership" "foo" {
  project_id = "${gitlab_project.foo.id}"
  user_id = "${gitlab_user.test.id}"
  expires_at = "2099-01-01"
  access_level = "guest"
}

resource "gitlab_project" "foo" {
  name = "foo%d"
  description = "Terraform acceptance tests"
 visibility_level ="public"
}

resource "gitlab_user" "test" {
  name = "foo%d"
  username = "listest%d"
  password = "test%dtt"
  email = "listest%d@ssss.com"
}
`, rInt, rInt, rInt, rInt, rInt)
}
