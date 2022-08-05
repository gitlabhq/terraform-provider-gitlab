//go:build acceptance
// +build acceptance

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccGitlabApplicationSettings_basic(t *testing.T) {
	// lintignore:AT001
	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			// Verify empty application settings
			{
				Config: `
					resource "gitlab_application_settings" "this" {}
				`,
			},
			// Verify changing some application settings
			{
				Config: `
					resource "gitlab_application_settings" "this" {
						after_sign_up_text = "Welcome to GitLab!"
					}
				`,
			},
		},
	})
}
