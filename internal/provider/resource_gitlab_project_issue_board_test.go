//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabProjectIssueBoard_basic(t *testing.T) {
	testProject := testAccCreateProject(t)
	testMilestone := testAccAddProjectMilestones(t, testProject, 1)[0]
	testLabels := testAccCreateProjectLabels(t, testProject.ID, 2)
	testUser := testAccCreateUsers(t, 1)[0]

	// NOTE: there is no way to delete the last issue board, see
	// https://gitlab.com/gitlab-org/gitlab/-/issues/367395
	testAccCreateProjectIssueBoard(t, testProject.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueBoardDestroy,
		Steps: []resource.TestStep{
			// Verify creation
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project = "%d"
						name    = "Test Board"
					}
				`, testProject.ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update with optional values (all optional attributes are EE only)
			{
				SkipFunc: isRunningInCE,
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"
						milestone_id = %d
						assignee_id  = %d
						labels       = ["%s", "%s"]
						weight       = 8
					}
				`, testProject.ID, testMilestone.ID, testUser.ID, testLabels[0].Name, testLabels[1].Name),
			},
			// Verify Import
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectIssueBoard_AllOnCreateEE(t *testing.T) {
	testAccCheckEE(t)

	testProject := testAccCreateProject(t)
	testMilestones := testAccAddProjectMilestones(t, testProject, 2)
	testLabels := testAccCreateProjectLabels(t, testProject.ID, 4)
	testUsers := testAccCreateUsers(t, 2)

	// NOTE: there is no way to delete the last issue board, see
	// https://gitlab.com/gitlab-org/gitlab/-/issues/367395
	testAccCreateProjectIssueBoard(t, testProject.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueBoardDestroy,
		Steps: []resource.TestStep{
			// Verify creation with all attributes set (some are only available in the update API)
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"
						milestone_id = %d
						assignee_id  = %d
						labels       = ["%s", "%s"]
						weight       = 8
					}
				`, testProject.ID, testMilestones[0].ID, testUsers[0].ID, testLabels[0].Name, testLabels[1].Name),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update with changed attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"
						milestone_id = %d
						assignee_id  = %d
						labels       = ["%s", "%s"]
						weight       = 8
					}
				`, testProject.ID, testMilestones[1].ID, testUsers[1].ID, testLabels[2].Name, testLabels[3].Name),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update with removed optional attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"
					}
				`, testProject.ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectIssueBoard_Lists(t *testing.T) {
	testProject := testAccCreateProject(t)
	testMilestones := testAccAddProjectMilestones(t, testProject, 2)
	testLabels := testAccCreateProjectLabels(t, testProject.ID, 4)
	testUsers := testAccCreateUsers(t, 2)
	testAccAddProjectMembers(t, testProject.ID, testUsers)

	// NOTE: there is no way to delete the last issue board, see
	// https://gitlab.com/gitlab-org/gitlab/-/issues/367395
	testAccCreateProjectIssueBoard(t, testProject.ID)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectIssueBoardDestroy,
		Steps: []resource.TestStep{
			// Create Board with 2 lists with core features
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"

						lists {
							label_id = %d
						}

						lists {
							label_id = %d
						}
					}
				`, testProject.ID, testLabels[0].ID, testLabels[1].ID),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update Board list labels
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"

						lists {
							label_id = %d
						}

						lists {
							label_id = %d
						}
					}
				`, testProject.ID, testLabels[2].ID, testLabels[3].ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Force a destroy for the board so that it can be recreated as the same resource
			{
				SkipFunc: isRunningInCE,
				Config:   ` `, // requires a space for empty config
			},
			{
				SkipFunc: isRunningInCE,
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"

						lists {
							label_id = %d
						}

						lists {
							assignee_id = %d
						}

						lists {
							milestone_id = %d
						}
					}
				`, testProject.ID, testLabels[0].ID, testUsers[0].ID, testMilestones[0].ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				SkipFunc: isRunningInCE,
				Config: fmt.Sprintf(`
					resource "gitlab_project_issue_board" "this" {
						project      = "%d"
						name         = "Test Board"

						lists {
							label_id = %d
						}

						lists {
							assignee_id = %d
						}

						lists {
							milestone_id = %d
						}
					}
				`, testProject.ID, testLabels[1].ID, testUsers[1].ID, testMilestones[1].ID),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_issue_board.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectIssueBoardDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_issue_board" {
			continue
		}

		project, issueBoardID, err := resourceGitlabProjectIssueBoardParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject, _, err := testGitlabClient.Boards.GetIssueBoard(project, issueBoardID)
		if err == nil && subject != nil {
			return fmt.Errorf("gitlab_project_issue_board resource '%s' still exists", rs.Primary.ID)
		}

		if err != nil && !is404(err) {
			return err
		}

		return nil
	}
	return nil
}
