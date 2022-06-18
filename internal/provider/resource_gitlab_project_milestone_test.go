//go:build acceptance
// +build acceptance

package provider

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabProjectMilestone_basic(t *testing.T) {

	rInt1, rInt2 := acctest.RandInt(), acctest.RandInt()
	project := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectMilestoneDestroy,
		Steps: []resource.TestStep{
			{
				// create Milestone with required values only
				Config: fmt.Sprintf(`
				resource "gitlab_project_milestone" "this" {
					project = "%v"
					title   = "test-%d"
				}`, project.PathWithNamespace, rInt1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "iid"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "milestone_id"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "updated_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "created_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "web_url"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "expired"),
				),
			},
			{
				// verify import
				ResourceName:      "gitlab_project_milestone.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				// update some Milestone attributes
				Config: fmt.Sprintf(`
				resource "gitlab_project_milestone" "this" {
					project     = "%[1]d"
					title       = "test-%[2]d"
					description = "test-%[2]d"
					start_date  = "2022-04-10"
					due_date    = "2022-04-15"
					state       = "closed"
				}`, project.ID, rInt2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "iid"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "milestone_id"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "updated_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "created_at"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "web_url"),
					resource.TestCheckResourceAttrSet("gitlab_project_milestone.this", "expired"),
				),
			},
			{
				// verify import
				ResourceName:      "gitlab_project_milestone.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectMilestoneDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_milestone" {
			continue
		}
		projectID, milestoneID, err := resourceGitLabProjectMilestoneParseId(rs.Primary.ID)
		if err != nil {
			return err
		}

		milestone, _, err := testGitlabClient.Milestones.GetMilestone(projectID, milestoneID)
		if err == nil && milestone != nil {
			return errors.New("Milestone still exists")
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}
