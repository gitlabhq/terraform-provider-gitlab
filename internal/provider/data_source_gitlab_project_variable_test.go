//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGitlabProjectVariable_basic(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_variable" "this" {
						project           = %d
						key               = "any_key"
					        value             = "any-value"
						environment_scope = "*"
					}

					data "gitlab_project_variable" "this" {
						project           = gitlab_project_variable.this.project
						key               = gitlab_project_variable.this.key
						environment_scope = gitlab_project_variable.this.environment_scope
					}
					`, testProject.ID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabProjectVariable("gitlab_project_variable.this", "data.gitlab_project_variable.this"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabProjectVariable(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resource := s.RootModule().Resources[src]
		resourceAttributes := resource.Primary.Attributes

		datasource := s.RootModule().Resources[n]
		datasourceAttributes := datasource.Primary.Attributes

		testAttributes := attributeNamesFromSchema(gitlabProjectVariableGetSchema())

		for _, attribute := range testAttributes {
			if datasourceAttributes[attribute] != resourceAttributes[attribute] {
				return fmt.Errorf("Expected variable's attribute `%s` to be: %s, but got: `%s`", attribute, resourceAttributes[attribute], datasourceAttributes[attribute])
			}
		}

		return nil
	}
}
