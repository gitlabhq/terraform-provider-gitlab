//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	. "github.com/onsi/gomega"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitLabProjectApprovalRule_Basic(t *testing.T) {
	// Set up project, groups, users, and branches to use in the test.

	testAccCheckEE(t)

	// Need to get the current user (usually the admin) because they are automatically added as group members, and we
	// will need the user ID for our assertions later.
	currentUser := testAccCurrentUser(t)

	project := testAccCreateProject(t)
	projectUsers := testAccCreateUsers(t, 2)
	branches := testAccCreateProtectedBranches(t, project, 2)
	groups := testAccCreateGroups(t, 2)
	group0Users := testAccCreateUsers(t, 1)
	group1Users := testAccCreateUsers(t, 1)

	testAccAddProjectMembers(t, project.ID, projectUsers) // Users must belong to the project for rules to work.
	testAccAddGroupMembers(t, groups[0].ID, group0Users)
	testAccAddGroupMembers(t, groups[1].ID, group1Users)

	// Terraform test starts here.

	var projectApprovalRule gitlab.ProjectApprovalRule

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectApprovalRuleDestroy(project.ID),
		Steps: []resource.TestStep{
			// Create rule
			{
				Config: testAccGitlabProjectApprovalRuleConfig_Basic(project.ID, 3, projectUsers[0].ID, groups[0].ID, branches[0].ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.foo", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes_Basic(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes_Basic{
						Name:                "foo",
						ApprovalsRequired:   3,
						EligibleApproverIDs: []int{currentUser.ID, projectUsers[0].ID, group0Users[0].ID},
						GroupIDs:            []int{groups[0].ID},
						ProtectedBranchIDs:  []int{branches[0].ID},
					}),
				),
			},
			// Update rule
			{
				Config: testAccGitlabProjectApprovalRuleConfig_Basic(project.ID, 2, projectUsers[1].ID, groups[1].ID, branches[1].ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.foo", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes_Basic(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes_Basic{
						Name:                "foo",
						ApprovalsRequired:   2,
						EligibleApproverIDs: []int{currentUser.ID, projectUsers[1].ID, group1Users[0].ID},
						GroupIDs:            []int{groups[1].ID},
						ProtectedBranchIDs:  []int{branches[1].ID},
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_approval_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitLabProjectApprovalRule_AnyApprover(t *testing.T) {
	// Set up project, groups, users, and branches to use in the test.

	testAccCheckEE(t)

	project := testAccCreateProject(t)

	// Terraform test starts here.

	var projectApprovalRule gitlab.ProjectApprovalRule

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectApprovalRuleDestroy(project.ID),
		Steps: []resource.TestStep{
			// Create rule
			{
				Config: testAccGitlabProjectApprovalRuleConfig_AnyApprover(project.ID, 3, "any_approver"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.bar", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes_AnyApprover(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes_AnyApprover{
						Name:              "bar",
						ApprovalsRequired: 3,
						RuleType:          "any_approver",
					}),
				),
			},
			// Update rule
			{
				Config: testAccGitlabProjectApprovalRuleConfig_AnyApprover(project.ID, 2, "any_approver"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.bar", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes_AnyApprover(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes_AnyApprover{
						Name:              "bar",
						ApprovalsRequired: 2,
						RuleType:          "any_approver",
					}),
				),
			},
			// Re-create rule
			{
				Config: testAccGitlabProjectApprovalRuleConfig_AnyApprover(project.ID, 2, "regular"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.bar", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes_AnyApprover(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes_AnyApprover{
						Name:              "bar",
						ApprovalsRequired: 2,
						RuleType:          "regular",
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_approval_rule.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccGitlabProjectApprovalRuleExpectedAttributes_Basic struct {
	Name                string
	ApprovalsRequired   int
	EligibleApproverIDs []int
	GroupIDs            []int
	ProtectedBranchIDs  []int
}

type testAccGitlabProjectApprovalRuleExpectedAttributes_AnyApprover struct {
	Name              string
	ApprovalsRequired int
	RuleType          string
}

func testAccCheckGitlabProjectApprovalRuleAttributes_Basic(got *gitlab.ProjectApprovalRule, want *testAccGitlabProjectApprovalRuleExpectedAttributes_Basic) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return InterceptGomegaFailure(func() {
			Expect(got.Name).To(Equal(want.Name), "name")
			Expect(got.ApprovalsRequired).To(Equal(want.ApprovalsRequired), "approvals_required")

			var approverIDs []int
			for _, approver := range got.EligibleApprovers {
				approverIDs = append(approverIDs, approver.ID)
			}
			Expect(approverIDs).To(ConsistOf(want.EligibleApproverIDs), "eligible_approvers")

			var groupIDs []int
			for _, group := range got.Groups {
				groupIDs = append(groupIDs, group.ID)
			}
			Expect(groupIDs).To(ConsistOf(want.GroupIDs), "groups")

			var protectedBranchIDs []int
			for _, branch := range got.ProtectedBranches {
				protectedBranchIDs = append(protectedBranchIDs, branch.ID)
			}
			Expect(protectedBranchIDs).To(ConsistOf(want.ProtectedBranchIDs), "protected_branches")
		})
	}
}

func testAccCheckGitlabProjectApprovalRuleAttributes_AnyApprover(got *gitlab.ProjectApprovalRule, want *testAccGitlabProjectApprovalRuleExpectedAttributes_AnyApprover) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return InterceptGomegaFailure(func() {
			Expect(got.Name).To(Equal(want.Name), "name")
			Expect(got.ApprovalsRequired).To(Equal(want.ApprovalsRequired), "approvals_required")
			Expect(got.RuleType).To(Equal(want.RuleType), "rule_type")
		})
	}
}

func testAccGitlabProjectApprovalRuleConfig_Basic(project, approvals, userID, groupID, protectedBranchID int) string {
	return fmt.Sprintf(`
resource "gitlab_project_approval_rule" "foo" {
  project              = %d
  name                 = "foo"
  approvals_required   = %d
  user_ids             = [%d]
  group_ids            = [%d]
  protected_branch_ids = [%d]
}`, project, approvals, userID, groupID, protectedBranchID)
}

func testAccGitlabProjectApprovalRuleConfig_AnyApprover(project, approvals int, rule_type string) string {
	return fmt.Sprintf(`
resource "gitlab_project_approval_rule" "bar" {
  project              = %d
  name                 = "bar"
  approvals_required   = %d
  rule_type            = "%s"
}`, project, approvals, rule_type)
}

func testAccCheckGitlabProjectApprovalRuleExists(n string, projectApprovalRule *gitlab.ProjectApprovalRule) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		projectID, ruleID, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return err
		}

		ruleIDInt, err := strconv.Atoi(ruleID)
		if err != nil {
			return err
		}

		rules, _, err := testGitlabClient.Projects.GetProjectApprovalRules(projectID)
		if err != nil {
			return err
		}

		for _, gotRule := range rules {
			if gotRule.ID == ruleIDInt {
				*projectApprovalRule = *gotRule
				return nil
			}
		}

		return fmt.Errorf("rule %d not found", ruleIDInt)
	}
}

func testAccCheckGitlabProjectApprovalRuleDestroy(pid interface{}) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		return InterceptGomegaFailure(func() {
			rules, _, err := testGitlabClient.Projects.GetProjectApprovalRules(pid)
			Expect(err).To(BeNil())
			Expect(rules).To(BeEmpty())
		})
	}
}
