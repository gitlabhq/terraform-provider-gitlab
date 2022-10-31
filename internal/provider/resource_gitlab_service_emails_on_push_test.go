//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceEmailsOnPush_basic(t *testing.T) {
	testProject := testAccCreateProject(t)

	var emailsOnPushService gitlab.EmailsOnPushService

	var recipients1 = "mynumberonerecipient@example.com"
	var recipients2 = "mynumbertworecipient@example.com"
	var emailsOnPushResourceName = "gitlab_service_emails_on_push.this"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabServiceEmailsOnPushDestroy,
		Steps: []resource.TestStep{
			// Create an Emails on Push service
			{
				Config: fmt.Sprintf(`
				resource "gitlab_service_emails_on_push" "this" {
					project    = %[1]d
					recipients = "%[2]s"
				}
				`, testProject.ID, recipients1),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceEmailsOnPushExists(emailsOnPushResourceName, &emailsOnPushService),
					resource.TestCheckResourceAttr(emailsOnPushResourceName, "recipients", recipients1),
					resource.TestCheckResourceAttr(emailsOnPushResourceName, "active", "true"),
					resource.TestCheckResourceAttrWith(emailsOnPushResourceName, "created_at", func(value string) error {
						expectedValue := emailsOnPushService.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_service_emails_on_push.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Emails on Push service
			{
				Config: fmt.Sprintf(`
				resource "gitlab_service_emails_on_push" "this" {
					project    = %[1]d
					recipients = "%[2]s"
				}
				`, testProject.ID, recipients2),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceEmailsOnPushExists(emailsOnPushResourceName, &emailsOnPushService),
					resource.TestCheckResourceAttr(emailsOnPushResourceName, "recipients", recipients2),
					resource.TestCheckResourceAttrWith(emailsOnPushResourceName, "created_at", func(value string) error {
						expectedValue := emailsOnPushService.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith(emailsOnPushResourceName, "updated_at", func(value string) error {
						expectedValue := emailsOnPushService.UpdatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_service_emails_on_push.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Emails on Push service to get back to previous settings
			{
				Config: fmt.Sprintf(`
				resource "gitlab_service_emails_on_push" "this" {
					project    = %[1]d
					recipients = "%[2]s"
				}
				`, testProject.ID, recipients1),
			},
			// Verify import
			{
				ResourceName:      "gitlab_service_emails_on_push.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabServiceEmailsOnPushExists(resourceIdentifier string, service *gitlab.EmailsOnPushService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceIdentifier]
		if !ok {
			return fmt.Errorf("Not Found: %s", resourceIdentifier)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}

		emailsOnPushService, _, err := testGitlabClient.Services.GetEmailsOnPushService(project)
		if err != nil {
			return fmt.Errorf("Emails on Push service does not exist in project %s: %v", project, err)
		}
		*service = *emailsOnPushService

		return nil
	}
}

func testAccCheckGitlabServiceEmailsOnPushDestroy(s *terraform.State) error {
	var project string

	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_service_emails_on_push" {
			project = rs.Primary.ID

			emailsOnPushService, _, err := testGitlabClient.Services.GetEmailsOnPushService(project)
			if err == nil {
				if emailsOnPushService != nil && emailsOnPushService.Active != false {
					return fmt.Errorf("[ERROR] Emails on Push Service %v still exists", project)
				}
			} else {
				return err
			}
		}
	}
	return nil
}
