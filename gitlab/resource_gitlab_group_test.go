package gitlab

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
	"net/http"
	"testing"
	"time"
)

func TestAccGitlabGroup_basic(t *testing.T) {
	var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Create a group
			{
				Config: testAccGitlabGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("foo-name-%d", rInt),
						Path:        fmt.Sprintf("foo-path-%d", rInt),
						Description: "Terraform acceptance tests",
						LFSEnabled:  true,
					}),
				),
			},
			// Update the group to change the description
			{
				Config: testAccGitlabGroupUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                 fmt.Sprintf("bar-name-%d", rInt),
						Path:                 fmt.Sprintf("bar-path-%d", rInt),
						Description:          "Terraform acceptance tests! Updated description",
						RequestAccessEnabled: true,
					}),
				),
			},
			// Update the group to put the name and description back
			{
				Config: testAccGitlabGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("foo-name-%d", rInt),
						Path:        fmt.Sprintf("foo-path-%d", rInt),
						Description: "Terraform acceptance tests",
						LFSEnabled:  true,
					}),
				),
			},
		},
	})
}

func TestAccGitlabGroup_subgroup(t *testing.T) {
	var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Create a group with subgroup_creation_level set to owner
			{
				Config: testAccGitlabSubGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                  fmt.Sprintf("foo-name-%d", rInt),
						Path:                  fmt.Sprintf("foo-path-%d", rInt),
						Description:           "Terraform acceptance tests",
						SubGroupCreationLevel: "owner",
						LFSEnabled:            true,
					}),
				),
			},
			// Update the group to change the level to maintainer
			{
				Config: testAccGitlabSubGroupUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                  fmt.Sprintf("foo-name-%d", rInt),
						Path:                  fmt.Sprintf("foo-path-%d", rInt),
						Description:           "Terraform acceptance tests",
						SubGroupCreationLevel: "maintainer",
						LFSEnabled:            true,
					}),
				),
			},
		},
	})
}

func TestAccGitlabGroup_projects(t *testing.T) {
	var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			// Create a group with project_creation_level set to noone
			{
				Config: testAccGitlabGroupProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                 fmt.Sprintf("foo-name-%d", rInt),
						Path:                 fmt.Sprintf("foo-path-%d", rInt),
						Description:          "Terraform acceptance tests",
						ProjectCreationLevel: "noone",
						LFSEnabled:           true,
					}),
				),
			},
			// Update the group to change the level to developer
			{
				Config: testAccGitlabGroupUpdateProjectConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                 fmt.Sprintf("foo-name-%d", rInt),
						Path:                 fmt.Sprintf("foo-path-%d", rInt),
						Description:          "Terraform acceptance tests",
						ProjectCreationLevel: "developer",
						LFSEnabled:           true,
					}),
				),
			},
			// Update the group to change the level to maintainer
			{
				Config: testAccGitlabGroupUpdateProjectConfigWithMaintainer(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupAttributes(&group, &testAccGitlabGroupExpectedAttributes{
						Name:                 fmt.Sprintf("foo-name-%d", rInt),
						Path:                 fmt.Sprintf("foo-path-%d", rInt),
						Description:          "Terraform acceptance tests",
						ProjectCreationLevel: "maintainer",
						LFSEnabled:           true,
					}),
				),
			},
		},
	})
}

func TestAccGitlabGroup_import(t *testing.T) {
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupConfig(rInt),
			},
			{
				ResourceName:      "gitlab_group.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabGroup_nested(t *testing.T) {
	var group gitlab.Group
	var group2 gitlab.Group
	var nestedGroup gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabNestedGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupExists("gitlab_group.foo2", &group2),
					testAccCheckGitlabGroupExists("gitlab_group.nested_foo", &nestedGroup),
					testAccCheckGitlabGroupAttributes(&nestedGroup, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("nfoo-name-%d", rInt),
						Path:        fmt.Sprintf("nfoo-path-%d", rInt),
						Description: "Terraform acceptance tests",
						LFSEnabled:  true,
						Parent:      &group,
					}),
				),
			},
			{
				Config: testAccGitlabNestedGroupChangeParentConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupExists("gitlab_group.foo2", &group2),
					testAccCheckGitlabGroupExists("gitlab_group.nested_foo", &nestedGroup),
					testAccCheckGitlabGroupAttributes(&nestedGroup, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("nfoo-name-%d", rInt),
						Path:        fmt.Sprintf("nfoo-path-%d", rInt),
						Description: "Terraform acceptance tests - new parent",
						LFSEnabled:  true,
						Parent:      &group2,
					}),
				),
			},
			{
				Config: testAccGitlabNestedGroupRemoveParentConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupExists("gitlab_group.foo2", &group2),
					testAccCheckGitlabGroupExists("gitlab_group.nested_foo", &nestedGroup),
					testAccCheckGitlabGroupAttributes(&nestedGroup, &testAccGitlabGroupExpectedAttributes{
						Name:        fmt.Sprintf("nfoo-name-%d", rInt),
						Path:        fmt.Sprintf("nfoo-path-%d", rInt),
						Description: "Terraform acceptance tests - updated",
						LFSEnabled:  true,
					}),
				),
			},
			// TODO In EE version, re-creating on the same path where a previous group was soft-deleted doesn't work.
			// {
			// 	Config: testAccGitlabNestedGroupConfig(rInt),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
			// 		testAccCheckGitlabGroupExists("gitlab_group.foo2", &group2),
			// 		testAccCheckGitlabGroupExists("gitlab_group.nested_foo", &nestedGroup),
			// 		testAccCheckGitlabGroupAttributes(&nestedGroup, &testAccGitlabGroupExpectedAttributes{
			// 			Name:        fmt.Sprintf("nfoo-name-%d", rInt),
			// 			Path:        fmt.Sprintf("nfoo-path-%d", rInt),
			// 			Description: "Terraform acceptance tests",
			// 			LFSEnabled:  true,
			// 			Parent:      &group,
			// 		}),
			// 	),
			// },
		},
	})
}

func TestAccGitlabGroup_disappears(t *testing.T) {
	var group gitlab.Group
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabGroupDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabGroupConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabGroupExists("gitlab_group.foo", &group),
					testAccCheckGitlabGroupDisappears(&group),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckGitlabGroupDisappears(group *gitlab.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := testAccProvider.Meta().(*gitlab.Client)

		_, err := conn.Groups.DeleteGroup(group.ID)
		if err != nil {
			return err
		}
		// Fixes groups API async deletion issue
		// https://github.com/terraform-providers/terraform-provider-gitlab/issues/319
		for start := time.Now(); time.Since(start) < 15*time.Second; {
			g, resp, err := conn.Groups.GetGroup(group.ID)
			if resp != nil && resp.StatusCode == http.StatusNotFound {
				return nil
			}
			if g != nil && g.MarkedForDeletionOn != nil {
				return nil
			}
			if err != nil {
				return err
			}
		}
		return fmt.Errorf("waited for more than 15 seconds for group to be asynchronously deleted")
	}
}

func testAccCheckGitlabGroupExists(n string, group *gitlab.Group) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		groupID := rs.Primary.ID
		if groupID == "" {
			return fmt.Errorf("No group ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotGroup, _, err := conn.Groups.GetGroup(groupID)
		if err != nil {
			return err
		}
		*group = *gotGroup
		return nil
	}
}

type testAccGitlabGroupExpectedAttributes struct {
	Name                  string
	Path                  string
	Description           string
	Parent                *gitlab.Group
	LFSEnabled            bool
	RequestAccessEnabled  bool
	ProjectCreationLevel  gitlab.ProjectCreationLevelValue
	SubGroupCreationLevel gitlab.SubGroupCreationLevelValue
}

func testAccCheckGitlabGroupAttributes(group *gitlab.Group, want *testAccGitlabGroupExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if group.Name != want.Name {
			return fmt.Errorf("got repo %q; want %q", group.Name, want.Name)
		}

		if group.Path != want.Path {
			return fmt.Errorf("got path %q; want %q", group.Path, want.Path)
		}

		if group.Description != want.Description {
			return fmt.Errorf("got description %q; want %q", group.Description, want.Description)
		}

		if group.LFSEnabled != want.LFSEnabled {
			return fmt.Errorf("got lfs_enabled %t; want %t", group.LFSEnabled, want.LFSEnabled)
		}

		if group.RequestAccessEnabled != want.RequestAccessEnabled {
			return fmt.Errorf("got request_access_enabled %t; want %t", group.RequestAccessEnabled, want.RequestAccessEnabled)
		}

		if want.Parent != nil {
			if group.ParentID != want.Parent.ID {
				return fmt.Errorf("got parent_id %d; want %d", group.ParentID, want.Parent.ID)
			}
		} else {
			if group.ParentID != 0 {
				return fmt.Errorf("got parent_id %d; want %d", group.ParentID, 0)
			}
		}

		if want.ProjectCreationLevel != "" {
			if group.ProjectCreationLevel != want.ProjectCreationLevel {
				return fmt.Errorf("got project_creation_level %q; want %q", group.ProjectCreationLevel, want.ProjectCreationLevel)
			}
		}

		if want.SubGroupCreationLevel != "" {
			if group.SubGroupCreationLevel != want.SubGroupCreationLevel {
				return fmt.Errorf("got subgroup_creation_level %q; want %q", group.SubGroupCreationLevel, want.SubGroupCreationLevel)
			}
		}

		return nil
	}
}

func testAccCheckGitlabGroupDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_group" {
			continue
		}

		group, resp, err := conn.Groups.GetGroup(rs.Primary.ID)
		if err == nil {
			if group != nil && fmt.Sprintf("%d", group.ID) == rs.Primary.ID {
				if group.MarkedForDeletionOn == nil {
					return fmt.Errorf("Group still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabGroupProjectConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"
  project_creation_level = "noone"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabGroupUpdateProjectConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"
  project_creation_level = "developer"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabGroupUpdateProjectConfigWithMaintainer(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"
  project_creation_level = "maintainer"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabSubGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"
  subgroup_creation_level = "owner"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabSubGroupUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"
  subgroup_creation_level = "maintainer"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabGroupUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "bar-name-%d"
  path = "bar-path-%d"
  description = "Terraform acceptance tests! Updated description"
  lfs_enabled = false
  request_access_enabled = true

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt)
}

func testAccGitlabNestedGroupConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "foo2" {
  name = "foo2-name-%d"
  path = "foo2-path-%d"
  description = "Terraform acceptance tests - parent2"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "nested_foo" {
  name = "nfoo-name-%d"
  path = "nfoo-path-%d"
  parent_id = "${gitlab_group.foo.id}"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabNestedGroupRemoveParentConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "foo2" {
  name = "foo2-name-%d"
  path = "foo2-path-%d"
  description = "Terraform acceptance tests - parent2"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "nested_foo" {
  name = "nfoo-name-%d"
  path = "nfoo-path-%d"
  description = "Terraform acceptance tests - updated"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt, rInt, rInt, rInt, rInt)
}

func testAccGitlabNestedGroupChangeParentConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_group" "foo" {
  name = "foo-name-%d"
  path = "foo-path-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "foo2" {
  name = "foo2-name-%d"
  path = "foo2-path-%d"
  description = "Terraform acceptance tests - parent2"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
resource "gitlab_group" "nested_foo" {
  name = "nfoo-name-%d"
  path = "nfoo-path-%d"
  description = "Terraform acceptance tests - new parent"
  parent_id = "${gitlab_group.foo2.id}"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}
  `, rInt, rInt, rInt, rInt, rInt, rInt)
}
