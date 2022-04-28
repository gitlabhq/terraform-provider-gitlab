package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataGitlabProjectMilestones_basic(t *testing.T) {
	testAccCheck(t)

	testProject := testAccCreateProject(t)
	testMilestones := testAccAddProjectMilestones(t, testProject, 2)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectMilestonesConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_project_milestones.this", "milestones.#", fmt.Sprintf("%d", len(testMilestones))),
					resource.TestCheckResourceAttr("data.gitlab_project_milestones.this", "milestones.0.title", testMilestones[1].Title),
					resource.TestCheckResourceAttr("data.gitlab_project_milestones.this", "milestones.0.description", testMilestones[1].Description),
					resource.TestCheckResourceAttr("data.gitlab_project_milestones.this", "milestones.1.title", testMilestones[0].Title),
					resource.TestCheckResourceAttr("data.gitlab_project_milestones.this", "milestones.1.description", testMilestones[0].Description),
				),
			},
		},
	})
}

func testAccDataGitlabProjectMilestonesConfig(projectID int) string {
	return fmt.Sprintf(`
data "gitlab_project_milestones" "this" {
  project_id   = "%d"
}
`, projectID)
}
