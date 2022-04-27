package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataGitlabProjectMilestones_basic(t *testing.T) {
	testAccCheck(t)
	countMilestones := 2
	project := testAccCreateProject(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectMilestones(countMilestones, project.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjectMilestones("gitlab_project_milestone.this", "data.gitlab_project_milestones.this", countMilestones),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectMilestones(src string, n string, countMilestones int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testAttributes := attributeNamesFromSchema(gitlabProjectMilestoneGetSchema())

		for numberMilestone := range make([]int, countMilestones) {
			search := s.RootModule().Resources[n]
			searchAttrs := search.Primary.Attributes

			milestone := s.RootModule().Resources[fmt.Sprintf("%s.%d", src, numberMilestone)]
			milestoneAttrs := milestone.Primary.Attributes

			for _, attribute := range testAttributes {
				milestoneAttr := milestoneAttrs[attribute]
				searchAttr := searchAttrs[fmt.Sprintf("milestones.%d.%s", numberMilestone, attribute)]
				if searchAttr != milestoneAttr {
					return fmt.Errorf("Expected the milestone `%s` with parameter `%s` to be: `%s`, but got: `%s`", milestoneAttrs["id"], attribute, milestoneAttr, searchAttr)
				}
			}
		}

		return nil
	}
}

func testAccDataGitlabProjectMilestones(countTags int, project int) string {
	return fmt.Sprintf(`
%s
data "gitlab_project_milestones" "this" {
  project_id = "%d"
  state      = "active"
  search     = "test"
  depends_on = [
    gitlab_project_milestone.this,
  ]
}
`, testAccDataGitlabProjectMilestonesSetup(countTags, project), project)
}

func testAccDataGitlabProjectMilestonesSetup(countTags int, project int) string {
	return fmt.Sprintf(`
resource "gitlab_project_milestone" "this" {
  count   = "%d"

  project_id  = "%d"
  title       = "test-${count.index}"
  description = "test-${count.index}"
}
`, countTags, project)
}
