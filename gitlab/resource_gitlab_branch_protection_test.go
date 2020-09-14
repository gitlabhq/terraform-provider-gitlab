package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabBranchProtection_basic(t *testing.T) {

	var pb gitlab.ProtectedBranch
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabBranchProtectionDestroy,
		Steps: []resource.TestStep{
			// Create a project and Branch Protection with default options
			{
				Config: testAccGitlabBranchProtectionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchProtectionExists("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionPersistsInStateCorrectly("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionAttributes(&pb, &testAccGitlabBranchProtectionExpectedAttributes{
						Name:                      fmt.Sprintf("BranchProtect-%d", rInt),
						PushAccessLevel:           accessLevel[gitlab.DeveloperPermissions],
						MergeAccessLevel:          accessLevel[gitlab.DeveloperPermissions],
						UnprotectAccessLevel:      accessLevel[gitlab.MaintainerPermissions],
						CodeOwnerApprovalRequired: false,
					}),
				),
			},
			// Update the Branch Protection
			{
				Config: testAccGitlabBranchProtectionUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchProtectionExists("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionPersistsInStateCorrectly("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionAttributes(&pb, &testAccGitlabBranchProtectionExpectedAttributes{
						Name:                      fmt.Sprintf("BranchProtect-%d", rInt),
						PushAccessLevel:           accessLevel[gitlab.MaintainerPermissions],
						MergeAccessLevel:          accessLevel[gitlab.MaintainerPermissions],
						UnprotectAccessLevel:      accessLevel[gitlab.DeveloperPermissions],
						CodeOwnerApprovalRequired: false,
					}),
				),
			},
			// Update code_owner_approval_required only (test Update instead of Create)
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabBranchProtectionUpdateCodeOwnerApprovalRequiredConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchProtectionExists("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionPersistsInStateCorrectly("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionAttributes(&pb, &testAccGitlabBranchProtectionExpectedAttributes{
						Name:                      fmt.Sprintf("BranchProtect-%d", rInt),
						PushAccessLevel:           accessLevel[gitlab.MaintainerPermissions],
						MergeAccessLevel:          accessLevel[gitlab.MaintainerPermissions],
						UnprotectAccessLevel:      accessLevel[gitlab.DeveloperPermissions],
						CodeOwnerApprovalRequired: true,
					}),
				),
			},
			// Update the Branch Protection using the more granular attributes
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabBranchProtectionUpdateGranularConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchProtectionExists("gitlab_branch_protection.branch_protect", &pb),
					resource.TestCheckResourceAttr("gitlab_branch_protection.branch_protect", "allowed_to_push.0.access_level.0", "maintainer"),
					resource.TestCheckResourceAttr("gitlab_branch_protection.branch_protect", "allowed_to_merge.0.access_level.0", "maintainer"),
					resource.TestCheckResourceAttr("gitlab_branch_protection.branch_protect", "allowed_to_merge.0.access_level.1", "developer"),
					resource.TestCheckResourceAttr("gitlab_branch_protection.branch_protect", "allowed_to_unprotect.0.access_level.0", "developer"),
				),
			},
			// Update the Branch Protection to get back to initial settings
			{
				Config: testAccGitlabBranchProtectionConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabBranchProtectionExists("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionPersistsInStateCorrectly("gitlab_branch_protection.branch_protect", &pb),
					testAccCheckGitlabBranchProtectionAttributes(&pb, &testAccGitlabBranchProtectionExpectedAttributes{
						Name:                      fmt.Sprintf("BranchProtect-%d", rInt),
						PushAccessLevel:           accessLevel[gitlab.DeveloperPermissions],
						MergeAccessLevel:          accessLevel[gitlab.DeveloperPermissions],
						UnprotectAccessLevel:      accessLevel[gitlab.MaintainerPermissions],
						CodeOwnerApprovalRequired: false,
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabBranchProtectionPersistsInStateCorrectly(n string, pb *gitlab.ProtectedBranch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		if rs.Primary.Attributes["merge_access_level"] != accessLevel[pb.MergeAccessLevels[0].AccessLevel] {
			return fmt.Errorf("merge access level not persisted in state correctly")
		}

		if rs.Primary.Attributes["push_access_level"] != accessLevel[pb.PushAccessLevels[0].AccessLevel] {
			return fmt.Errorf("push access level not persisted in state correctly")
		}

		return nil
	}
}

func testAccCheckGitlabBranchProtectionExists(n string, pb *gitlab.ProtectedBranch) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}
		project, branch, err := projectAndBranchFromID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in Splitting Project and Branch Ids")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		pbs, _, err := conn.ProtectedBranches.ListProtectedBranches(project, nil)
		if err != nil {
			return err
		}
		for _, gotpb := range pbs {
			if gotpb.Name == branch {
				*pb = *gotpb
				return nil
			}
		}
		return fmt.Errorf("Protected Branch does not exist")
	}
}

type testAccGitlabBranchProtectionExpectedAttributes struct {
	Name                      string
	PushAccessLevel           string
	MergeAccessLevel          string
	UnprotectAccessLevel      string
	CodeOwnerApprovalRequired bool
}

func testAccCheckGitlabBranchProtectionAttributes(pb *gitlab.ProtectedBranch, want *testAccGitlabBranchProtectionExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pb.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", pb.Name, want.Name)
		}

		if pb.PushAccessLevels[0].AccessLevel != accessLevelID[want.PushAccessLevel] {
			return fmt.Errorf("got Push access levels %q; want %q", pb.PushAccessLevels[0].AccessLevel, accessLevelID[want.PushAccessLevel])
		}

		if pb.MergeAccessLevels[0].AccessLevel != accessLevelID[want.MergeAccessLevel] {
			return fmt.Errorf("got Merge access levels %q; want %q", pb.MergeAccessLevels[0].AccessLevel, accessLevelID[want.MergeAccessLevel])
		}

		if pb.UnprotectAccessLevels != nil {
			if pb.UnprotectAccessLevels[0].AccessLevel != accessLevelID[want.UnprotectAccessLevel] {
				return fmt.Errorf("got Unprotect access levels %q; want %q", pb.UnprotectAccessLevels[0].AccessLevel, accessLevelID[want.UnprotectAccessLevel])
			}
		}

		if pb.CodeOwnerApprovalRequired != want.CodeOwnerApprovalRequired {
			return fmt.Errorf("got code owner approval required %t; want %t", pb.CodeOwnerApprovalRequired, want.CodeOwnerApprovalRequired)
		}

		return nil
	}
}

func testAccCheckGitlabBranchProtectionDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)
	var project string
	var branch string
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_branch_protection" {
			branch = rs.Primary.ID
		}
	}

	pb, response, err := conn.ProtectedBranches.GetProtectedBranch(project, branch)
	if err == nil {
		if pb != nil {
			return fmt.Errorf("project branch protection %s still exists", branch)
		}
	}
	if response.StatusCode != 404 {
		return err
	}
	return nil
}

func testAccGitlabBranchProtectionConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_branch_protection" "branch_protect" {
  project = gitlab_project.foo.id
  branch = "BranchProtect-%d"
  push_access_level = "developer"
  merge_access_level = "developer"
  # Omit unprotect_access_level to test default value set by GitLab
}
	`, rInt, rInt)
}

func testAccGitlabBranchProtectionUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_branch_protection" "branch_protect" {
	project = gitlab_project.foo.id
	branch = "BranchProtect-%d"
	push_access_level = "maintainer"
	merge_access_level = "maintainer"
	unprotect_access_level = "developer"
}
	`, rInt, rInt)
}

func testAccGitlabBranchProtectionUpdateCodeOwnerApprovalRequiredConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_branch_protection" "branch_protect" {
	project = gitlab_project.foo.id
	branch = "BranchProtect-%d"
	push_access_level = "maintainer"
	merge_access_level = "maintainer"
	unprotect_access_level = "developer"
	code_owner_approval_required = true
}
	`, rInt, rInt)
}

func testAccGitlabBranchProtectionUpdateGranularConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_branch_protection" "branch_protect" {
	project = gitlab_project.foo.id
	branch = "BranchProtect-%d"
	allowed_to_push {
		access_level = ["maintainer"]
	}
	allowed_to_merge {
        access_level = ["maintainer", "developer"]
    }
	allowed_to_unprotect {
		access_level = ["developer"]
	}
}
	`, rInt, rInt)
}
