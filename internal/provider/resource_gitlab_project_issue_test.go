//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectIssue_basic(t *testing.T) {
	var testIssue gitlab.Issue
	var updatedTestIssue gitlab.Issue

	testProject := testAccCreateProject(t)
	testUser := testAccCreateUsers(t, 1)[0]
	testAccAddProjectMembers(t, testProject.ID, []*gitlab.User{testUser})
	testMilestone := testAccAddProjectMilestones(t, testProject, 1)[0]

	currentUser, _, err := testGitlabClient.Users.CurrentUser()
	if err != nil {
		t.Fatalf("Failed to get current user: %v", err)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueDestroy,
		Steps: []resource.TestStep{
			// create Issue with required values only
			{
				Config: testAccGitlabProjectIssueConfigRequiredOnly(testProject),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectIssueExists("gitlab_project_issue.this", &testIssue),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "project", testProject.PathWithNamespace),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "iid", "1"),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "title", "Terraform test issue"),
					resource.TestCheckResourceAttrWith("gitlab_project_issue.this", "created_at", func(value string) error {
						expectedValue := testIssue.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
			// update some Issue attributes
			{
				Config: testAccGitlabProjectIssueConfigAll(testProject, testMilestone, testUser),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectIssueExists("gitlab_project_issue.this", &updatedTestIssue),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "project", testProject.PathWithNamespace),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "iid", "1"),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "title", "Terraform test issue"),
					resource.TestCheckResourceAttrWith("gitlab_project_issue.this", "updated_at", func(value string) error {
						expectedValue := updatedTestIssue.UpdatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
			// go back to required values only
			{
				Config: testAccGitlabProjectIssueConfigRequiredOnly(testProject),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectIssueExists("gitlab_project_issue.this", &testIssue),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "project", testProject.PathWithNamespace),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "iid", "1"),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "title", "Terraform test issue"),
					func(s *terraform.State) error {
						if !updatedTestIssue.UpdatedAt.Before(*testIssue.UpdatedAt) {
							return fmt.Errorf("expected issue to be updated, but it wasn't")
						}
						return nil
					},
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
			// close issue
			{
				Config: testAccGitlabProjectIssueConfigWithState(testProject, "closed"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectIssueExists("gitlab_project_issue.this", &testIssue),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "state", "closed"),
					resource.TestCheckResourceAttrWith("gitlab_project_issue.this", "closed_at", func(value string) error {
						expectedValue := testIssue.ClosedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "closed_by_user_id", fmt.Sprintf("%d", currentUser.ID)),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
			// re-open issue
			{
				Config: testAccGitlabProjectIssueConfigWithState(testProject, "opened"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectIssueExists("gitlab_project_issue.this", &testIssue),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "state", "opened"),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "closed_at", ""),
					resource.TestCheckResourceAttr("gitlab_project_issue.this", "closed_by_user_id", "0"),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
		},
	})
}

func TestAccGitlabProjectIssue_basicEE(t *testing.T) {
	testAccCheckEE(t)

	testProject := testAccCreateProject(t)
	testUser := testAccCreateUsers(t, 1)[0]
	testAccAddProjectMembers(t, testProject.ID, []*gitlab.User{testUser})
	testMilestone := testAccAddProjectMilestones(t, testProject, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueDestroy,
		Steps: []resource.TestStep{
			// create Issue with EE features set
			{
				Config: testAccGitlabProjectIssueConfigEE(testProject, testMilestone, testUser),
			},
			// Verify import
			{
				ResourceName:            "gitlab_project_issue.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"delete_on_destroy"},
			},
		},
	})
}

func TestAccGitlabProjectIssue_deleteOnDestroy(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueDestroy,
		Steps: []resource.TestStep{
			// create Issue with required values only
			{
				Config: testAccGitlabProjectIssueConfigDeleteOnDestroy(testProject),
			},
		},
	})
}

func testAccCheckGitlabProjectIssueExists(n string, issue *gitlab.Issue) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		project, issueIID, err := resourceGitLabProjectIssueParseId(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error parsing issue ID: %s", err)
		}

		gotIssue, _, err := testGitlabClient.Issues.GetIssue(project, issueIID)
		if err != nil {
			return fmt.Errorf("Cannot get issue: %v", err)
		}

		log.Printf("[DEBUG] testAccCheckGitlabProjectIssueExists: %#v", gotIssue)

		*issue = *gotIssue
		return nil
	}
}

func testAccCheckGitlabProjectIssueDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_issue" {
			continue
		}

		project, issueIID, err := resourceGitLabProjectIssueParseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		issue, _, err := testGitlabClient.Issues.GetIssue(project, issueIID)
		if err == nil && issue != nil && issue.IID == issueIID {
			if val, ok := rs.Primary.Attributes["delete_on_destory"]; ok && val == "true" {
				return fmt.Errorf("Issue still exists")
			} else {
				if issue.State != "closed" {
					return fmt.Errorf("Issue still in state %s (should be closed)", issue.State)
				}
			}
		}
		if err != nil && !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabProjectIssueConfigRequiredOnly(project *gitlab.Project) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
	// required
	project = "%s"
	title   = "Terraform test issue"

	// we have to add this because of: https://github.com/xanzy/go-gitlab/issues/1384
	// optional
	due_date = "2020-01-01"
}`, project.PathWithNamespace)
}

func testAccGitlabProjectIssueConfigWithState(project *gitlab.Project, state string) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
	project  = "%s"
	title    = "Terraform test issue"
	// we have to add this because of: https://github.com/xanzy/go-gitlab/issues/1384
	due_date = "2020-01-01"

	state    = "%s"
}`, project.PathWithNamespace, state)
}

func testAccGitlabProjectIssueConfigDeleteOnDestroy(project *gitlab.Project) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
	project = "%s"
	title   = "Terraform test issue"

	delete_on_destroy = true
}`, project.PathWithNamespace)
}

func testAccGitlabProjectIssueConfigAll(project *gitlab.Project, milestone *gitlab.Milestone, assignee *gitlab.User) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
	// required
	project = "%s"
	title   = "Terraform test issue"

	// we have to add this because of: https://github.com/xanzy/go-gitlab/issues/1384
	// optional
	due_date = "2020-01-01"

	assignee_ids = [%d]
	confidential = true
	description  = "Terraform test issue description"
	issue_type   = "issue"
	labels       = ["foo", "bar"]
	milestone_id = %d
	state 	     = "opened"
	discussion_locked = true
}`, project.PathWithNamespace, assignee.ID, milestone.ID)
}

func testAccGitlabProjectIssueConfigEE(project *gitlab.Project, milestone *gitlab.Milestone, assignee *gitlab.User) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
	// required
	project = "%s"
	title   = "Terraform test issue"

	// we have to add this because of: https://github.com/xanzy/go-gitlab/issues/1384
	// optional
	due_date = "2020-01-01"

	assignee_ids = [%d]
	confidential = true
	description  = "Terraform test issue description"
	issue_type   = "issue"
	labels       = ["foo", "bar"]
	milestone_id = %d
	weight       = 42
	state 	     = "opened"
	discussion_locked = true
}`, project.PathWithNamespace, assignee.ID, milestone.ID)
}
