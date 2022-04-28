package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabProjectMilestone_basic(t *testing.T) {
	testAccCheck(t)

	testProject := testAccCreateProject(t)
	testMilestone := testAccAddProjectMilestones(t, testProject, 1)[0]

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectMilestoneConfig(testProject.ID, testMilestone.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "milestone_id", fmt.Sprintf("%v", testMilestone.ID)),
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "title", testMilestone.Title),
					resource.TestCheckResourceAttr("data.gitlab_project_milestone.this", "description", testMilestone.Description),
				),
			},
		},
	})
}

func testAccDataGitlabProjectMilestoneConfig(projectID int, milestoneID int) string {
	return fmt.Sprintf(`
data "gitlab_project_milestone" "this" {
  project_id   = "%d"
  milestone_id = "%d"
}
`, projectID, milestoneID)
}
