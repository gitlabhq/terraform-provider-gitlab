//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabProjectBadge_basic(t *testing.T) {
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectBadgeDestroy,
		Steps: []resource.TestStep{
			// Create a project badge
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_badge" "this" {
						project   = "%d"
						link_url  = "https://example.com/badge"
						image_url = "https://example.com/badge.svg"
						name      = "badge"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_project_badge.this", "rendered_link_url"),
					resource.TestCheckResourceAttrSet("gitlab_project_badge.this", "rendered_image_url"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_badge.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the project badge
			{
				Config: fmt.Sprintf(`
					resource "gitlab_project_badge" "this" {
						project   = "%d"
						link_url  = "https://example.com/badge-updated"
						image_url = "https://example.com/badge-updated.svg"
						name      = "badge-updated"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_project_badge.this", "rendered_link_url"),
					resource.TestCheckResourceAttrSet("gitlab_project_badge.this", "rendered_image_url"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_project_badge.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectBadgeDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project_badge" {
			continue
		}

		projectID, badgeID, err := resourceGitlabProjectBadgeParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		gotBadge, _, err := testGitlabClient.ProjectBadges.GetProjectBadge(projectID, badgeID)
		if err != nil {
			if !is404(err) {
				return err
			}
			return nil
		}

		if gotBadge != nil {
			return fmt.Errorf("Badge still exists")
		}

		return nil
	}
	return nil
}
