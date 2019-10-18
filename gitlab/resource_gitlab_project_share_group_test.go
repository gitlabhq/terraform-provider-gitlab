package gitlab

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectShareGroup_basic(t *testing.T) {
	var membership struct {
		GroupID          int    "json:\"group_id\""
		GroupName        string "json:\"group_name\""
		GroupAccessLevel int    "json:\"group_access_level\""
	}
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{PreCheck: func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectShareGroupDestroy,
		Steps: []resource.TestStep{

			// Assign member to the project as a developer
			{
				Config: testAccGitlabProjectShareGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(testAccCheckGitlabProjectShareGroupExists("gitlab_project_share_group.foo", &membership), testAccCheckGitlabProjectShareGroupAttributes(&membership, &testAccGitlabProjectShareGroupExpectedAttributes{
					access_level: fmt.Sprintf("developer"),
				})),
			},
		},
	})
}

func testAccCheckGitlabProjectShareGroupExists(n string, membership *struct {
	GroupID          int    "json:\"group_id\""
	GroupName        string "json:\"group_name\""
	GroupAccessLevel int    "json:\"group_access_level\""
}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		conn := testAccProvider.Meta().(*gitlab.Client)
		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		projectID := rs.Primary.Attributes["project_id"]
		if projectID == "" {
			return fmt.Errorf("no project ID is set")
		}

		groupID := rs.Primary.Attributes["group_id"]
		if groupID == "" {
			return fmt.Errorf("no user id is set")
		}

		gotProjectShareGroup, _, err := conn.Projects.GetProject(projectID, nil)
		if err != nil {
			return err
		}

		*membership = gotProjectShareGroup.SharedWithGroups[0]
		return nil
	}
}

type testAccGitlabProjectShareGroupExpectedAttributes struct {
	access_level string
}

func testAccCheckGitlabProjectShareGroupAttributes(membership *struct {
	GroupID          int    "json:\"group_id\""
	GroupName        string "json:\"group_name\""
	GroupAccessLevel int    "json:\"group_access_level\""
}, want *testAccGitlabProjectShareGroupExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		accessLevelId, ok := accessLevel[gitlab.AccessLevelValue(membership.GroupAccessLevel)]
		if !ok {
			return fmt.Errorf("invalid access level '%s'", accessLevelId)
		}
		if accessLevelId != want.access_level {
			return fmt.Errorf("got access level %s; want %s", accessLevelId, want.access_level)
		}
		return nil
	}
}

func testAccCheckGitlabProjectShareGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_share_group" {
			continue
		}

		projectID := rs.Primary.Attributes["project_id"]
		groupID := rs.Primary.Attributes["group_id"]

		// GetProjectMember needs int type for groupID
		groupIDI, err := strconv.Atoi(groupID)
		gotShareGroup, _, err := conn.Projects.GetProject(projectID, nil)
		if err != nil {
			return nil
		}

		for _, v := range gotShareGroup.SharedWithGroups {
			if groupIDI == v.GroupID {
				if gotShareGroup != nil && fmt.Sprintf("%d", v.GroupAccessLevel) == rs.Primary.Attributes["access_level"] {
					return fmt.Errorf("project still has a reference.")
				}
			}
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectShareGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_share_group" "foo" {
  project_id = "${gitlab_project.foo.id}"
  group_id = "${gitlab_group.test.id}"
  access_level = "developer"
}

resource "gitlab_project" "foo" {
  name = "foo%d"
  description = "Terraform acceptance tests"
  visibility_level ="public"

  lifecycle {
    ignore_changes = [
   		shared_with_groups
    ]
  }
}

resource "gitlab_group" "test" {
  name        = "foo%d"
  path        = "foo%d"
  description = "Description for foo%d"
}

`, rInt, rInt, rInt, rInt)
}
