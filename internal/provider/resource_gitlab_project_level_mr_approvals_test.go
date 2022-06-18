//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectLevelMRApprovals_basic(t *testing.T) {
	testAccCheckEE(t)

	var projectApprovals gitlab.ProjectApprovals
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectLevelMRApprovalsDestroy,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_level_mr_approvals" "foo" {
						project_id                                     = "%d"
						reset_approvals_on_push                        = true
						disable_overriding_approvers_per_merge_request = true
						merge_requests_author_approval                 = true
						merge_requests_disable_committers_approval     = true
						require_password_to_approve                    = true
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      true,
						disableOverridingApproversPerMergeRequest: true,
						mergeRequestsAuthorApproval:               true,
						mergeRequestsDisableCommittersApproval:    true,
						requirePasswordToApprove:                  true,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_level_mr_approvals.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_level_mr_approvals" "foo" {
						project_id                                     = "%d"
						reset_approvals_on_push                        = false
						disable_overriding_approvers_per_merge_request = false
						merge_requests_author_approval                 = false
						merge_requests_disable_committers_approval     = false
						require_password_to_approve                    = false
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      false,
						disableOverridingApproversPerMergeRequest: false,
						mergeRequestsAuthorApproval:               false,
						mergeRequestsDisableCommittersApproval:    false,
						requirePasswordToApprove:                  false,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_level_mr_approvals.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_level_mr_approvals" "foo" {
						project_id                                     = "%d"
						reset_approvals_on_push                        = true
						disable_overriding_approvers_per_merge_request = true
						merge_requests_author_approval                 = true
						merge_requests_disable_committers_approval     = true
						require_password_to_approve                    = true
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectLevelMRApprovalsExists("gitlab_project_level_mr_approvals.foo", &projectApprovals),
					testAccCheckGitlabProjectLevelMRApprovalsAttributes(&projectApprovals, &testAccGitlabProjectLevelMRApprovalsExpectedAttributes{
						resetApprovalsOnPush:                      true,
						disableOverridingApproversPerMergeRequest: true,
						mergeRequestsAuthorApproval:               true,
						mergeRequestsDisableCommittersApproval:    true,
						requirePasswordToApprove:                  true,
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_level_mr_approvals.foo",
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
	requirePasswordToApprove                  bool
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
		if projectApprovals.RequirePasswordToApprove != want.requirePasswordToApprove {
			return fmt.Errorf("got require_password_to_approve %t; want %t", projectApprovals.RequirePasswordToApprove, want.requirePasswordToApprove)
		}
		return nil
	}
}

func testAccCheckGitlabProjectLevelMRApprovalsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if resp != nil && resp.StatusCode != 404 {
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

		gotApprovalConfig, _, err := testGitlabClient.Projects.GetApprovalConfiguration(projectId)
		if err != nil {
			return err
		}

		*projectApprovals = *gotApprovalConfig
		return nil
	}
}
