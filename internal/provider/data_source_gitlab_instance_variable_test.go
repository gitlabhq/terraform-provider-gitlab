//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccDataSourceGitlabInstanceVariable_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
					resource "gitlab_instance_variable" "this" {
						key               = "any_key"
					        value             = "any-value"
					}

					data "gitlab_instance_variable" "this" {
						key               = gitlab_instance_variable.this.key
					}
				`,
				Check: resource.ComposeTestCheckFunc(
					testAccDataSourceGitlabInstanceVariable("gitlab_instance_variable.this", "data.gitlab_instance_variable.this"),
				),
			},
		},
	})
}

func testAccDataSourceGitlabInstanceVariable(src, n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		resource := s.RootModule().Resources[src]
		resourceAttributes := resource.Primary.Attributes

		datasource := s.RootModule().Resources[n]
		datasourceAttributes := datasource.Primary.Attributes

		testAttributes := attributeNamesFromSchema(gitlabInstanceVariableGetSchema())

		for _, attribute := range testAttributes {
			if datasourceAttributes[attribute] != resourceAttributes[attribute] {
				return fmt.Errorf("Expected variable's attribute `%s` to be: %s, but got: `%s`", attribute, resourceAttributes[attribute], datasourceAttributes[attribute])
			}
		}

		return nil
	}
}
