package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectApprovalsConfiguration_basic(t *testing.T) {

	var approvals gitlab.ProjectApprovals
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabProjectApprovalsConfigurationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabProjectApprovalsConfigurationUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalsConfigurationExists("gitlab_project_approvals_configuration.approvals", &approvals),
					testAccCheckGitlabProjectApprovalsConfigurationAttributes(&approvals, &testAccGitlabProjectApprovalsConfigurationExpectedAttributes{
						approvalsBeforeMerge:                      2,
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

func testAccGitlabProjectApprovalsConfigurationUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
	name              = "foo-%d"
	description       = "Terraform acceptance tests"
	visibility_level  = "public"
}

resource "gitlab_project_approvals_configuration" "approvals" {
	project                                        = gitlab_project.foo.id
	approvals_before_merge						   = 2
	reset_approvals_on_push                        = true
	disable_overriding_approvers_per_merge_request = true
	merge_requests_author_approval                 = true
	merge_requests_disable_committers_approval     = true
}
	`, rInt)
}

func testAccCheckGitlabProjectApprovalsConfigurationDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			continue
		}
		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				return fmt.Errorf("Repository still exists")
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func testAccCheckGitlabProjectApprovalsConfigurationExists(n string, approvals *gitlab.ProjectApprovals) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.ID
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		gotApprovals, _, err := conn.Projects.GetApprovalConfiguration(project)
		if err != nil {
			return err
		}
		*approvals = *gotApprovals
		return nil
	}
}

type testAccGitlabProjectApprovalsConfigurationExpectedAttributes struct {
	approvalsBeforeMerge                      int
	resetApprovalsOnPush                      bool
	disableOverridingApproversPerMergeRequest bool
	mergeRequestsAuthorApproval               bool
	mergeRequestsDisableCommittersApproval    bool
}

func testAccCheckGitlabProjectApprovalsConfigurationAttributes(approvals *gitlab.ProjectApprovals, want *testAccGitlabProjectApprovalsConfigurationExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if approvals.ApprovalsBeforeMerge != want.approvalsBeforeMerge {
			return fmt.Errorf("got description %q; want %q", approvals.ApprovalsBeforeMerge, want.approvalsBeforeMerge)
		}
		if approvals.ResetApprovalsOnPush != want.resetApprovalsOnPush {
			return fmt.Errorf("got reset approvals on push %t; want %t", approvals.ResetApprovalsOnPush, want.resetApprovalsOnPush)
		}
		if approvals.DisableOverridingApproversPerMergeRequest != want.disableOverridingApproversPerMergeRequest {
			return fmt.Errorf("got disable overriding approvers per merge request %t; want %t", approvals.DisableOverridingApproversPerMergeRequest, want.disableOverridingApproversPerMergeRequest)
		}
		if approvals.MergeRequestsAuthorApproval != want.mergeRequestsAuthorApproval {
			return fmt.Errorf("got allow merge request author approval %t; want %t", approvals.MergeRequestsAuthorApproval, want.mergeRequestsAuthorApproval)
		}
		if approvals.MergeRequestsDisableCommittersApproval != want.mergeRequestsDisableCommittersApproval {
			return fmt.Errorf("got disable merge request commiters approval %t; want %t", approvals.MergeRequestsDisableCommittersApproval, want.mergeRequestsDisableCommittersApproval)
		}
		return nil
	}
}
