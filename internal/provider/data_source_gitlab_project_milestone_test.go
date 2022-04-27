package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGitlabProjectMilestone_basic(t *testing.T) {
	testAccCheck(t)

	testProject := testAccCreateProject(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectMilestoneConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjectMilestone("gitlab_project_milestone.this", "data.gitlab_project_milestone.this"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectMilestone(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resource := s.RootModule().Resources[src]
		resourceAttributes := resource.Primary.Attributes

		datasource := s.RootModule().Resources[n]
		datasourceAttributes := datasource.Primary.Attributes

		testAttributes := attributeNamesFromSchema(gitlabProjectMilestoneGetSchema())

		for _, attribute := range testAttributes {
			if datasourceAttributes[attribute] != resourceAttributes[attribute] {
				return fmt.Errorf("Expected issue's attribute `%s` to be: %s, but got: `%s`", attribute, resourceAttributes[attribute], datasourceAttributes[attribute])
			}
		}

		return nil
	}
}

func testAccDataGitlabProjectMilestoneConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_project_milestone" "this" {
  project_id  = %[1]d
  title       = "Terraform acceptance tests"
  description = "Some description"
  start_date  = "1994-02-21"
  due_date    = "1994-02-25"
}
data "gitlab_project_milestone" "this" {
  project_id   = %[1]d
  milestone_id = gitlab_project_milestone.this.milestone_id
}
`, projectID)
}
