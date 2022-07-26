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

func TestAccGitlabServiceJira_basic(t *testing.T) {
	var jiraService gitlab.JiraService
	rInt := acctest.RandInt()
	jiraResourceName := "gitlab_service_jira.jira"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabServiceJiraDestroy,
		Steps: []resource.TestStep{
			// Create a project and a jira service
			{
				Config: testAccGitlabServiceJiraConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJiraExists(jiraResourceName, &jiraService),
					resource.TestCheckResourceAttr(jiraResourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(jiraResourceName, "username", "user1"),
					resource.TestCheckResourceAttr(jiraResourceName, "password", "mypass"),
					resource.TestCheckResourceAttr(jiraResourceName, "commit_events", "true"),
					resource.TestCheckResourceAttr(jiraResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(jiraResourceName, "comment_on_event_enabled", "false"),
				),
			},
			// Verify Import
			{
				ResourceName:      jiraResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: as soon as we remove support for GitLab < 15.2 we can remove ignoring `jira_issue_transition_id`.
				//        See https://gitlab.com/gitlab-org/gitlab/-/issues/362437
				ImportStateVerifyIgnore: []string{"password", "jira_issue_transition_id"},
			},
			// Update the jira service
			{
				Config: testAccGitlabServiceJiraUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJiraExists(jiraResourceName, &jiraService),
					resource.TestCheckResourceAttr(jiraResourceName, "url", "https://testurl.com"),
					resource.TestCheckResourceAttr(jiraResourceName, "api_url", "https://testurl.com/rest"),
					resource.TestCheckResourceAttr(jiraResourceName, "username", "user2"),
					resource.TestCheckResourceAttr(jiraResourceName, "password", "mypass_update"),
					resource.TestCheckResourceAttr(jiraResourceName, "jira_issue_transition_id", "3"),
					resource.TestCheckResourceAttr(jiraResourceName, "commit_events", "false"),
					resource.TestCheckResourceAttr(jiraResourceName, "merge_requests_events", "true"),
					resource.TestCheckResourceAttr(jiraResourceName, "comment_on_event_enabled", "true"),
				),
			},
			// Verify Import
			{
				ResourceName:      jiraResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: as soon as we remove support for GitLab < 15.2 we can remove ignoring `jira_issue_transition_id`.
				//        See https://gitlab.com/gitlab-org/gitlab/-/issues/362437
				ImportStateVerifyIgnore: []string{"password", "jira_issue_transition_id"},
			},
			// Update the jira service to get back to previous settings
			{
				Config: testAccGitlabServiceJiraConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceJiraExists(jiraResourceName, &jiraService),
					resource.TestCheckResourceAttr(jiraResourceName, "url", "https://test.com"),
					resource.TestCheckResourceAttr(jiraResourceName, "api_url", "https://testurl.com/rest"),
					resource.TestCheckResourceAttr(jiraResourceName, "username", "user1"),
					resource.TestCheckResourceAttr(jiraResourceName, "password", "mypass"),
					resource.TestCheckResourceAttr(jiraResourceName, "commit_events", "true"),
					resource.TestCheckResourceAttr(jiraResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(jiraResourceName, "comment_on_event_enabled", "false"),
				),
			},
			// Verify Import
			{
				ResourceName:      jiraResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				// TODO: as soon as we remove support for GitLab < 15.2 we can remove ignoring `jira_issue_transition_id`.
				//        See https://gitlab.com/gitlab-org/gitlab/-/issues/362437
				ImportStateVerifyIgnore: []string{"password", "jira_issue_transition_id"},
			},
		},
	})
}

func testAccCheckGitlabServiceJiraExists(n string, service *gitlab.JiraService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		jiraService, _, err := testGitlabClient.Services.GetJiraService(project)
		if err != nil {
			return fmt.Errorf("Jira service does not exist in project %s: %v", project, err)
		}
		*service = *jiraService

		return nil
	}
}

func testAccCheckGitlabServiceJiraDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_service_jira" {
			continue
		}

		project := rs.Primary.ID

		_, _, err := testGitlabClient.Services.GetJiraService(project)
		if err == nil {
			return fmt.Errorf("Jira Service Integration in project %s still exists", project)
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func testAccGitlabServiceJiraConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_jira" "jira" {
  project  = "${gitlab_project.foo.id}"
  url      = "https://test.com"
  username = "user1"
  password = "mypass"
  commit_events = true
  merge_requests_events    = false
  comment_on_event_enabled = false
}
`, rInt)
}

func testAccGitlabServiceJiraUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_jira" "jira" {
  project  = "${gitlab_project.foo.id}"
  url      = "https://testurl.com"
  api_url  = "https://testurl.com/rest"
  username = "user2"
  password = "mypass_update"
  jira_issue_transition_id = "3"
  commit_events = false
  merge_requests_events    = true
  comment_on_event_enabled = true
}
`, rInt)
}
