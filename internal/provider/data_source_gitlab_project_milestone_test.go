//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabProjectMilestone_basic(t *testing.T) {

	testProject := testAccCreateProject(t)
	testMilestone := testAccAddProjectMilestones(t, testProject, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
				data "gitlab_project_milestone" "this" {
					project      = "%d"
					milestone_id = "%d"
				}`, testProject.ID, testMilestone.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "milestone_id", fmt.Sprintf("%v", testMilestone.ID)),
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "title", testMilestone.Title),
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "description", testMilestone.Description),
				),
			},
		},
	})
}
