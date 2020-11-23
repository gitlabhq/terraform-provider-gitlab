package gitlab

import (
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitLabProjectApprovalRule_basic(t *testing.T) {
	var projectApprovalRule gitlab.ProjectApprovalRule
	randomInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckGitlabProjectApprovalRuleDestroy,
		Steps: []resource.TestStep{
			{ // Create Rule
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectApprovalRuleConfig(randomInt, 3, "", "gitlab_group.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.foo", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes{
						ApproverUsernames: []string{fmt.Sprintf("foo-user-%d", randomInt)},
						ApprovalsRequired: 3,
						GroupPaths:        []string{fmt.Sprintf("foo-group-%d", randomInt)},
						Name:              fmt.Sprintf("foo rule %d", randomInt),
						RandomInt:         randomInt,
					}),
				),
			},
			{ // Add group and user
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectApprovalRuleConfig(randomInt, 2, "gitlab_user.baz.id", "gitlab_group.bar.id, gitlab_group.foo.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.foo", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes{
						ApproverUsernames: []string{
							fmt.Sprintf("bar-user-%d", randomInt),
							fmt.Sprintf("baz-user-%d", randomInt),
							fmt.Sprintf("foo-user-%d", randomInt),
						},
						ApprovalsRequired: 2,
						GroupPaths: []string{
							fmt.Sprintf("bar-group-%d", randomInt),
							fmt.Sprintf("foo-group-%d", randomInt),
						},
						Name:      fmt.Sprintf("foo rule %d", randomInt),
						RandomInt: randomInt,
					}),
				),
			},
			{ // Remove group and user
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectApprovalRuleConfig(randomInt, 1, "gitlab_user.qux.id", "gitlab_group.bar.id"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectApprovalRuleExists("gitlab_project_approval_rule.foo", &projectApprovalRule),
					testAccCheckGitlabProjectApprovalRuleAttributes(&projectApprovalRule, &testAccGitlabProjectApprovalRuleExpectedAttributes{
						ApproverUsernames: []string{
							fmt.Sprintf("bar-user-%d", randomInt),
							fmt.Sprintf("qux-user-%d", randomInt),
						},
						ApprovalsRequired: 1,
						GroupPaths:        []string{fmt.Sprintf("bar-group-%d", randomInt)},
						Name:              fmt.Sprintf("foo rule %d", randomInt),
						RandomInt:         randomInt,
					}),
				),
			},
		},
	})
}

func TestAccGitLabProjectApprovalRule_import(t *testing.T) {
	randomInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		PreCheck:     func() { testAccPreCheck(t) },
		CheckDestroy: testAccCheckGitlabProjectApprovalRuleDestroy,
		Steps: []resource.TestStep{
			{ // Create Rule
				SkipFunc: isRunningInCE,
				Config:   testAccGitLabProjectApprovalRuleConfig(randomInt, 1, "", "gitlab_group.foo.id"),
			},
			{ // Verify Import
				SkipFunc:          isRunningInCE,
				ResourceName:      "gitlab_project_approval_rule.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

type testAccGitlabProjectApprovalRuleExpectedAttributes struct {
	ApprovalsRequired int
	ApproverUsernames []string
	GroupPaths        []string
	Name              string
	RandomInt         int
}

func testAccCheckGitlabProjectApprovalRuleAttributes(projectApprovalRule *gitlab.ProjectApprovalRule, want *testAccGitlabProjectApprovalRuleExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if projectApprovalRule.Name != want.Name {
			return fmt.Errorf("got name %s; want %s", projectApprovalRule.Name, want.Name)
		}

		if projectApprovalRule.ApprovalsRequired != want.ApprovalsRequired {
			return fmt.Errorf("got approvals_required %d; want %d", projectApprovalRule.ApprovalsRequired, want.ApprovalsRequired)
		}

		// Compare unique usernames
		var userNames []string
		for _, approver := range projectApprovalRule.EligibleApprovers {
			// Approvers will include the group creator, which will come from the GITLAB_TOKEN user.
			// Filter for users with RandomInt in the username
			if strings.HasSuffix(approver.Username, strconv.Itoa(want.RandomInt)) {
				userNames = append(userNames, approver.Username)
			}
		}
		sort.Strings(userNames)

		if !reflect.DeepEqual(userNames, want.ApproverUsernames) {
			return fmt.Errorf("got approvers %s; want %s", userNames, want.ApproverUsernames)
		}

		// Compare unique group paths
		var groupPaths []string
		for _, group := range projectApprovalRule.Groups {
			groupPaths = append(groupPaths, group.Path)
		}
		sort.Strings(groupPaths)

		if !reflect.DeepEqual(groupPaths, want.GroupPaths) {
			return fmt.Errorf("got groups %s; want %s", groupPaths, want.GroupPaths)
		}

		return nil
	}
}

func testAccGitLabProjectApprovalRuleConfig(
	randomInt int,
	approvals int,
	userIDs string,
	groupIDs string,
) string {
	return fmt.Sprintf(`
resource "gitlab_user" "foo" {
	name             = "foo user"
	username         = "foo-user-%[1]d"
	password         = "foo12345"
	email            = "foo-user%[1]d@ssss.com"
	is_admin         = false
  projects_limit   = 2
  can_create_group = false
  is_external      = false
}

resource "gitlab_user" "bar" {
	name             = "bar user"
	username         = "bar-user-%[1]d"
	password         = "bar12345"
	email            = "bar-user%[1]d@ssss.com"
	is_admin         = false
  projects_limit   = 2
  can_create_group = false
  is_external      = false
}

resource "gitlab_user" "baz" {
	name             = "baz user"
	username         = "baz-user-%[1]d"
	password         = "baz12345"
	email            = "baz-user%[1]d@ssss.com"
	is_admin         = false
  projects_limit   = 2
  can_create_group = false
  is_external      = false
}

resource "gitlab_user" "qux" {
	name             = "qux user"
	username         = "qux-user-%[1]d"
	password         = "qux12345"
	email            = "qux-user%[1]d@ssss.com"
	is_admin         = false
  projects_limit   = 2
  can_create_group = false
  is_external      = false
}

resource "gitlab_project" "foo" {
	name              = "foo project %[1]d"
	path              = "foo-project-%[1]d"
	description       = "Terraform acceptance test - Approval Rule"
	visibility_level  = "public"
}

resource "gitlab_project_membership" "baz" {
  project_id     = gitlab_project.foo.id
  user_id        = gitlab_user.baz.id
  access_level   = "developer"
}

resource "gitlab_project_membership" "qux" {
  project_id     = gitlab_project.foo.id
  user_id        = gitlab_user.qux.id
  access_level   = "developer"
}

resource "gitlab_group" "foo" {
	name             = "foo-group %[1]d"
	path             = "foo-group-%[1]d"
	description      = "Terraform acceptance tests - Approval Rule"
	visibility_level = "public"
}

resource "gitlab_group" "bar" {
	name             = "bar-group %[1]d"
	path             = "bar-group-%[1]d"
	description      = "Terraform acceptance tests - Approval Rule"
	visibility_level = "public"
}

resource "gitlab_group_membership" "foo" {
  group_id         = gitlab_group.foo.id
  user_id          = gitlab_user.foo.id
  access_level     = "developer"
}

resource "gitlab_group_membership" "bar" {
  group_id        = gitlab_group.bar.id
  user_id         = gitlab_user.bar.id
  access_level    = "developer"
}

resource "gitlab_project_approval_rule" "foo" {
	project            = gitlab_project.foo.id
	name               = "foo rule %[1]d"
	approvals_required = %d
	user_ids           = [%s]
	group_ids          = [%s]
}
	`,
		randomInt,
		approvals,
		userIDs,
		groupIDs,
	)
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

		client := testAccProvider.Meta().(*gitlab.Client)

		rules, _, err := client.Projects.GetProjectApprovalRules(projectID)
		if err != nil {
			return err
		}

		for _, gotRule := range rules {
			if gotRule.ID == ruleIDInt {
				*projectApprovalRule = *gotRule
				return nil
			}
		}

		return fmt.Errorf("Rule %d not found", ruleIDInt)
	}
}

func testAccCheckGitlabProjectApprovalRuleDestroy(s *terraform.State) error {
	err := testAccCheckGitlabProjectDestroy(s)
	if err != nil {
		return err
	}

	err = testAccCheckGitlabGroupDestroy(s)
	if err != nil {
		return err
	}

	err = testAccCheckGitlabUserDestroy(s)
	if err != nil {
		return err
	}

	return nil
}
