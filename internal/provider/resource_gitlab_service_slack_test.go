//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceSlack_basic(t *testing.T) {
	var slackService gitlab.SlackService
	rInt := acctest.RandInt()
	slackResourceName := "gitlab_service_slack.slack"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabServiceSlackDestroy,
		Steps: []resource.TestStep{
			// Create a project and a slack service with minimal settings
			{
				Config: testAccGitlabServiceSlackMinimalConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
				),
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
			// Update slack service with more settings
			{
				Config: testAccGitlabServiceSlackConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", "test"),
					// TODO: Currently, GitLab doesn't correctly implement the API, so this is
					//       impossible to implement here at the moment.
					//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
					// resource.TestCheckResourceAttr(slackResourceName, "deployment_events", "true"),
					// resource.TestCheckResourceAttr(slackResourceName, "deployment_channel", "test"),
					resource.TestCheckResourceAttr(slackResourceName, "notify_only_broken_pipelines", "true"),
				),
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
			// Update the slack service to get back to minimal settings
			{
				Config: testAccGitlabServiceSlackMinimalConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceExists(slackResourceName, &slackService),
					resource.TestCheckResourceAttr(slackResourceName, "webhook", "https://test.com"),
					resource.TestCheckResourceAttr(slackResourceName, "push_channel", ""),
				),
			},
			// Verify Import
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
		slackService, _, err := testGitlabClient.Services.GetSlackService(project)
		if err != nil {
			return fmt.Errorf("Slack service does not exist in project %s: %v", project, err)
		}
		*service = *slackService

		return nil
	}
}

func testAccCheckGitlabServiceSlackDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_service_slack" {
			continue
		}

		project := rs.Primary.ID

		_, _, err := testGitlabClient.Services.GetSlackService(project)
		if err == nil {
			return fmt.Errorf("Slack Service Integration in project %s still exists", project)
		}
		if !is404(err) {
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
// TODO: Currently, GitLab doesn't correctly implement the API, so this is
//       impossible to implement here at the moment.
//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
//   deployment_channel           = "test"
//   deployment_events            = true
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
// TODO: Currently, GitLab doesn't correctly implement the API, so this is
//       impossible to implement here at the moment.
//       see https://gitlab.com/gitlab-org/gitlab/-/issues/28903
//   deployment_channel           = "test deployment_channel"
//   deployment_events            = false
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
  branches_to_be_notified      = "all"
}
`, rInt)
}
