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

func TestAccGitlabServiceMicrosoftTeams_basic(t *testing.T) {
	var teamsService gitlab.MicrosoftTeamsService
	rInt := acctest.RandInt()
	teamsResourceName := "gitlab_service_microsoft_teams.teams"

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabServiceMicrosoftTeamsDestroy,
		Steps: []resource.TestStep{
			// Create a project and a teams service
			{
				Config: testAccGitlabServiceMicrosoftTeamsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceMicrosoftTeamsExists(teamsResourceName, &teamsService),
					resource.TestCheckResourceAttr(teamsResourceName, "webhook", "https://test.com/?token=4"),
					resource.TestCheckResourceAttr(teamsResourceName, "notify_only_broken_pipelines", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "branches_to_be_notified", "all"),
					resource.TestCheckResourceAttr(teamsResourceName, "push_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "issues_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_issues_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "tag_push_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "note_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_note_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "pipeline_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "wiki_page_events", "false"),
				),
			},
			// Update the teams service
			{
				Config: testAccGitlabServiceMicrosoftTeamsUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceMicrosoftTeamsExists(teamsResourceName, &teamsService),
					resource.TestCheckResourceAttr(teamsResourceName, "webhook", "https://testurl.com/?token=5"),
					resource.TestCheckResourceAttr(teamsResourceName, "notify_only_broken_pipelines", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "branches_to_be_notified", "default"),
					resource.TestCheckResourceAttr(teamsResourceName, "push_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "issues_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_issues_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "merge_requests_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "tag_push_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "note_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_note_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "pipeline_events", "true"),
					resource.TestCheckResourceAttr(teamsResourceName, "wiki_page_events", "true"),
				),
			},
			// Update the teams service to get back to previous settings
			{
				Config: testAccGitlabServiceMicrosoftTeamsConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceMicrosoftTeamsExists(teamsResourceName, &teamsService),
					resource.TestCheckResourceAttr(teamsResourceName, "webhook", "https://test.com/?token=4"),
					resource.TestCheckResourceAttr(teamsResourceName, "notify_only_broken_pipelines", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "branches_to_be_notified", "all"),
					resource.TestCheckResourceAttr(teamsResourceName, "push_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "issues_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_issues_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "merge_requests_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "tag_push_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "note_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "confidential_note_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "pipeline_events", "false"),
					resource.TestCheckResourceAttr(teamsResourceName, "wiki_page_events", "false"),
				),
			},
			{
				ResourceName:      teamsResourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabServiceMicrosoftTeamsExists(n string, service *gitlab.MicrosoftTeamsService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		teamsService, _, err := testGitlabClient.Services.GetMicrosoftTeamsService(project)
		if err != nil {
			return fmt.Errorf("Microsoft Teams service does not exist in project %s: %v", project, err)
		}
		*service = *teamsService

		return nil
	}
}

func testAccCheckGitlabServiceMicrosoftTeamsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, resp, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
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

func testAccGitlabServiceMicrosoftTeamsConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_microsoft_teams" "teams" {
  project  = "${gitlab_project.foo.id}"
  webhook = "https://test.com/?token=4"
  notify_only_broken_pipelines = false
  branches_to_be_notified = "all"
  push_events = false
  issues_events = false
  confidential_issues_events = false
  merge_requests_events = false
  tag_push_events = false
  note_events = false
  confidential_note_events = false
  pipeline_events = false
  wiki_page_events = false
}
`, rInt)
}

func testAccGitlabServiceMicrosoftTeamsUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project" "foo" {
  name        = "foo-%d"
  description = "Terraform acceptance tests"
  # So that acceptance tests can be run in a gitlab organization
  # with no billing
  visibility_level = "public"
}

resource "gitlab_service_microsoft_teams" "teams" {
  project  = "${gitlab_project.foo.id}"
  webhook = "https://testurl.com/?token=5"
  notify_only_broken_pipelines = true
  branches_to_be_notified = "default"
  push_events = true
  issues_events = true
  confidential_issues_events = true
  merge_requests_events = true
  tag_push_events = true
  note_events = true
  confidential_note_events = true
  pipeline_events = true
  wiki_page_events = true
}
`, rInt)
}
