//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGitlabGroupVariable_basic(t *testing.T) {
	testGroup := testAccCreateGroups(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					resource "gitlab_group_variable" "this" {
						group             = %d
						key               = "any_key"
					        value             = "any-value"
						environment_scope = "*"
					}

					data "gitlab_group_variable" "this" {
						group             = gitlab_group_variable.this.group
						key               = gitlab_group_variable.this.key
						environment_scope = gitlab_group_variable.this.environment_scope
					}
					`, testGroup.ID,
				),
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabGroupVariable("gitlab_group_variable.this", "data.gitlab_group_variable.this"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabGroupVariable(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resource := s.RootModule().Resources[src]
		resourceAttributes := resource.Primary.Attributes

		datasource := s.RootModule().Resources[n]
		datasourceAttributes := datasource.Primary.Attributes

		testAttributes := attributeNamesFromSchema(gitlabGroupVariableGetSchema())

		for _, attribute := range testAttributes {
			if datasourceAttributes[attribute] != resourceAttributes[attribute] {
				return fmt.Errorf("Expected variable's attribute `%s` to be: %s, but got: `%s`", attribute, resourceAttributes[attribute], datasourceAttributes[attribute])
			}
		}

		return nil
	}
}
