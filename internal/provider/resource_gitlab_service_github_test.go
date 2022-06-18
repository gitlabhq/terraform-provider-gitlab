//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceGithub_basic(t *testing.T) {
	testAccCheckEE(t)

	var githubService gitlab.GithubService
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabServiceGithubDestroy,
		Steps: []resource.TestStep{
			// Create a project and a github service
			{
				Config: fmt.Sprintf(`
					resource "gitlab_service_github" "github" {
						project        = "%d"
						token          = "test"
						repository_url = "https://github.com/gitlabhq/terraform-provider-gitlab"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists("gitlab_service_github.github", &githubService),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "repository_url", "https://github.com/gitlabhq/terraform-provider-gitlab"),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "static_context", "true"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_service_github.github",
				ImportStateIdFunc: getGithubProjectID("gitlab_service_github.github"),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
				},
			},
			// Update the github service
			{
				Config: fmt.Sprintf(`
					resource "gitlab_service_github" "github" {
						project        = "%d"
						token          = "test"
						repository_url = "https://github.com/terraform-providers/terraform-provider-github"
						static_context = false
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists("gitlab_service_github.github", &githubService),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "repository_url", "https://github.com/terraform-providers/terraform-provider-github"),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "static_context", "false"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_service_github.github",
				ImportStateIdFunc: getGithubProjectID("gitlab_service_github.github"),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
				},
			},
			// Update the github service to get back to previous settings
			{
				Config: fmt.Sprintf(`
					resource "gitlab_service_github" "github" {
						project        = "%d"
						token          = "test"
						repository_url = "https://github.com/gitlabhq/terraform-provider-gitlab"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists("gitlab_service_github.github", &githubService),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "repository_url", "https://github.com/gitlabhq/terraform-provider-gitlab"),
					resource.TestCheckResourceAttr("gitlab_service_github.github", "static_context", "true"),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_service_github.github",
				ImportStateIdFunc: getGithubProjectID("gitlab_service_github.github"),
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"token",
				},
			},
		},
	})
}

func testAccCheckGitlabServiceGithubExists(n string, service *gitlab.GithubService) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}
		githubService, _, err := testGitlabClient.Services.GetGithubService(project)
		if err != nil {
			return fmt.Errorf("Github service does not exist in project %s: %v", project, err)
		}
		*service = *githubService

		return nil
	}
}

func testAccCheckGitlabServiceGithubDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_project" {
			continue
		}

		gotRepo, _, err := testGitlabClient.Projects.GetProject(rs.Primary.ID, nil)
		if err == nil {
			if gotRepo != nil && fmt.Sprintf("%d", gotRepo.ID) == rs.Primary.ID {
				if gotRepo.MarkedForDeletionAt == nil {
					return fmt.Errorf("Repository still exists")
				}
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func getGithubProjectID(n string) resource.ImportStateIdFunc {
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
