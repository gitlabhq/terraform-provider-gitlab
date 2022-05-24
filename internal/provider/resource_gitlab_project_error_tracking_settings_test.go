//go:build acceptance
// +build acceptance

package provider

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabProjectErrorTrackingSettings_basic(t *testing.T) {
	// Set up project.
	project := testAccCreateProject(t)

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectErrorTrackingSettingsDestroy(project.ID),
		Steps: []resource.TestStep{
			// Configure Error Tracking settings for a project.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_error_tracking_settings" "this" {
					project_id = %d
					integrated = true
				  }`, project.ID),
				ExpectError: regexp.MustCompile("Error Tracking must be enabled in the GitLab UI first before using this resource."),
			},
			// Unfortunately as it currently stands Error Tracking must be enabled in the UI before this resource can be used
			// The best we can do is catch the error that is thrown if this resource is attempted to be used without UI enablement first
		},
	})
}

func testAccCheckGitlabProjectErrorTrackingSettingsDestroy(projectID int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		ErrorTrackingSettings, _, err := testGitlabClient.ErrorTracking.GetErrorTrackingSettings(projectID)
		if err != nil {
			if is404(err) {
				return nil
			}
			return errors.New("Error disabling Error Tracking for Project")
		}

		if ErrorTrackingSettings.Active == true {
			return fmt.Errorf("Error tracking is still enabled for project: %w", err)
		}

		return nil
	}
}
