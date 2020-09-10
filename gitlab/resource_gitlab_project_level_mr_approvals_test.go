package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectLevelMRApprovals_basic(t *testing.T) {

	var projectApprovals gitlab.ProjectApprovals
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectLevelMRApprovalsDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectLevelMRApprovalsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      true,
						disableOverridingApproversPerMergeRequest: true,
						mergeRequestsAuthorApproval:               true,
						mergeRequestsDisableCommittersApproval:    true,
					}),
				),
			},
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectLevelMRApprovalsUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      false,
						disableOverridingApproversPerMergeRequest: false,
						mergeRequestsAuthorApproval:               false,
						mergeRequestsDisableCommittersApproval:    false,
					}),
				),
			},
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectLevelMRApprovalsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      true,
						disableOverridingApproversPerMergeRequest: true,
						mergeRequestsAuthorApproval:               true,
						mergeRequestsDisableCommittersApproval:    true,
					}),
				),
			},
		},
	})
}

func TestAccGitlabProjectLevelMRApprovals_import(t *testing.T) {
	resourceName := "gitlab_project_level_mr_approvals.foo"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectLevelMRApprovalsDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectLevelMRApprovalsConfig(rInt),
			},
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccGitlabProjectLevelMRApprovalsExpectedAttributes struct {
	resetApprovalsOnPush                      bool
	disableOverridingApproversPerMergeRequest bool
	mergeRequestsAuthorApproval               bool
	mergeRequestsDisableCommittersApproval    bool
}

func testAccCheckGitlabProjectLevelMRApprovalsAttributes(projectApprovals *gitlab.ProjectApprovals, want *testAccGitlabProjectLevelMRApprovalsExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if projectApprovals.ResetApprovalsOnPush != want.resetApprovalsOnPush {
			return fmt.Errorf("got reset_approvals_on_push %t; want %t", projectApprovals.ResetApprovalsOnPush, want.resetApprovalsOnPush)
		}
		if projectApprovals.DisableOverridingApproversPerMergeRequest != want.disableOverridingApproversPerMergeRequest {
			return fmt.Errorf("got disable_overriding_approvers_per_merge_request %t; want %t", projectApprovals.DisableOverridingApproversPerMergeRequest, want.disableOverridingApproversPerMergeRequest)
		}
		if projectApprovals.MergeRequestsAuthorApproval != want.mergeRequestsAuthorApproval {
			return fmt.Errorf("got merge_requests_author_approval %t; want %t", projectApprovals.MergeRequestsAuthorApproval, want.mergeRequestsAuthorApproval)
		}
		if projectApprovals.MergeRequestsDisableCommittersApproval != want.mergeRequestsDisableCommittersApproval {
			return fmt.Errorf("got merge_requests_disable_committers_approval %t; want %t", projectApprovals.MergeRequestsDisableCommittersApproval, want.mergeRequestsDisableCommittersApproval)
		}
		return nil
	}
}

func testAccCheckGitlabProjectLevelMRApprovalsDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists.")
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

func testAccCheckGitlabProjectLevelMRApprovalsExists(n string, projectApprovals *gitlab.ProjectApprovals) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		projectId := rs.Primary.ID
		if projectId == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotApprovalConfig, _, err := conn.Projects.GetApprovalConfiguration(projectId)
		if err != nil {
			return err
		}

		*projectApprovals = *gotApprovalConfig
		return nil
	}
}

func testAccGitlabProjectLevelMRApprovalsConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
	name              = "foo-%d"
	description       = "Terraform acceptance tests"
	visibility_level  = "public"
}

resource "gitlab_project_level_mr_approvals" "foo" {
	project_id                                     = gitlab_project.foo.id
	reset_approvals_on_push                        = true
	disable_overriding_approvers_per_merge_request = true
	merge_requests_author_approval                 = true
	merge_requests_disable_committers_approval     = true
}
	`, rInt)
}

func testAccGitlabProjectLevelMRApprovalsUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
	name              = "foo-%d"
	description       = "Terraform acceptance tests"
	visibility_level  = "public"
}

resource "gitlab_project_level_mr_approvals" "foo" {
	project_id                                     = gitlab_project.foo.id
	reset_approvals_on_push                        = false
	disable_overriding_approvers_per_merge_request = false
	merge_requests_author_approval                 = false
	merge_requests_disable_committers_approval     = false
}
	`, rInt)
}
