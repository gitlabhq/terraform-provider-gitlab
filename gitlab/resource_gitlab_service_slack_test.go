package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceSlack_basic(t *testing.T) {
	var slackService gitlab.SlackService
	rInt := acctest.RandInt()
	slackResourceName := "gitlab_service_slack.slack"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceSlackDestroy,
		Steps: []resource.TestStep{
			// Create a project and a slack service with minimal settings
			{
				Config: testAccGitlabServiceSlackMinimalConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
				),
			},
			// Update slack service with more settings
			{
				Config: testAccGitlabServiceSlackConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", "test"),
					resource.TestCheckResourceAttr(slackResourceName, "notify_only_broken_pipelines", "true"),
				),
			},
			// Update the slack service
			{
				Config: testAccGitlabServiceSlackUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://testwebhook.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_events", "false"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", "test push_channel"),
					resource.TestCheckResourceAttr(slackResourceName, "notify_only_broken_pipelines", "false"),
				),
			},
			// Update the slack service to get back to previous settings
			{
				Config: testAccGitlabServiceSlackConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", "test"),
					resource.TestCheckResourceAttr(slackResourceName, "notify_only_broken_pipelines", "true"),
				),
			},
			// Update the slack service to get back to minimal settings
			{
				Config: testAccGitlabServiceSlackMinimalConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", ""),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabServiceSlack_import(t *testing.T) {
	slackResourceName := "gitlab_service_slack.slack"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceSlackDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabServiceSlackConfig(rInt),
			},
			{
				ResourceName:      slackResourceName,
				ImportStateIdFunc: getSlackProjectID(slackResourceName),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"notify_only_broken_pipelines",
					"notify_only_default_branch",
				},
			},
		},
	})
}

func testAccCheckGitlabServiceExists(n string, service *gitlab.SlackService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		conn := testAccProvider.Meta().(*gitlab.Client)

		slackService, _, err := conn.Services.GetSlackService(project)
		if err != nil {
			return fmt.Errorf("Slack service does not exist in project %s: %v", project, err)
		}
		*service = *slackService

		return nil
	}
}

func testAccCheckGitlabServiceSlackDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*gitlab.Client)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := conn.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if resp.StatusCode != 404 {
			return err
		}
		return nil
	}
	return nil
}

func getSlackProjectID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return project, nil
	}
}

func testAccGitlabServiceSlackMinimalConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_slack" "slack" {
  project                      = "${gitlab_project.foo.id}"
  webhook                      = "https://test.com"
}
`, rInt)
}

func testAccGitlabServiceSlackConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_slack" "slack" {
  project                      = "${gitlab_project.foo.id}"
  webhook                      = "https://test.com"
  username                     = "test"
  push_events                  = true
  push_channel                 = "test"
  issues_events                = true
  issue_channel                = "test"
  confidential_issues_events   = true
  confidential_issue_channel   = "test"
  confidential_note_events     = true
  merge_requests_events        = true
  merge_request_channel        = "test"
  tag_push_events              = true
  tag_push_channel             = "test"
  note_events                  = true
  note_channel                 = "test"
  pipeline_events              = true
  pipeline_channel             = "test"
  wiki_page_events             = true
  wiki_page_channel            = "test"
  notify_only_broken_pipelines = true
  notify_only_default_branch   = true
  branches_to_be_notified      = "all"
}
`, rInt)
}

func testAccGitlabServiceSlackUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"

  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_slack" "slack" {
  project                      = "${gitlab_project.foo.id}"
  webhook                      = "https://testwebhook.com"
  username                     = "test username"
  push_events                  = false
  push_channel                 = "test push_channel"
  issues_events                = false
  issue_channel                = "test issue_channel"
  confidential_issues_events   = false
  confidential_issue_channel   = "test confidential_issue_channel"
  confidential_note_events     = false
  merge_requests_events        = false
  merge_request_channel        = "test merge_request_channel"
  tag_push_events              = false
  tag_push_channel             = "test tag_push_channel"
  note_events                  = false
  note_channel                 = "test note_channel"
  pipeline_events              = false
  pipeline_channel             = "test pipeline_channel"
  wiki_page_events             = false
  wiki_page_channel            = "test wiki_page_channel"
  notify_only_broken_pipelines = false
  notify_only_default_branch   = false
  branches_to_be_notified      = "all"
}
`, rInt)
}
