//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGitlabProjectIssue_basic(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccDataGitlabProjectIssueConfig(testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjectIssue("gitlab_project_issue.this", "data.gitlab_project_issue.this"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectIssue(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resource := s.RootModule().Resources[src]
		resourceAttributes := resource.Primary.Attributes

		datasource := s.RootModule().Resources[n]
		datasourceAttributes := datasource.Primary.Attributes

		testAttributes := attributeNamesFromSchema(gitlabProjectIssueGetSchema())

		for _, attribute := range testAttributes {
			if datasourceAttributes[attribute] != resourceAttributes[attribute] {
				return fmt.Errorf("Expected issue's attribute `%s` to be: %s, but got: `%s`", attribute, resourceAttributes[attribute], datasourceAttributes[attribute])
			}
		}

		return nil
	}
}

func testAccDataGitlabProjectIssueConfig(projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_project_issue" "this" {
  project     = %d
  title       = "Terraform acceptance tests"
  description = "Some description"
  due_date    = "1994-02-21"
}

data "gitlab_project_issue" "this" {
	project = %d
	iid     = gitlab_project_issue.this.iid
}
`, projectID, projectID)
}
