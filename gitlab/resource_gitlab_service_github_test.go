package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabServiceGithub_basic(t *testing.T) {
	var githubService gitlab.GithubService
	rInt := acctest.RandInt()
	githubResourceName := "gitlab_service_github.github"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceGithubDestroy,
		Steps: []resource.TestStep{
			// Create a project and a github service
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabServiceGithubConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists(githubResourceName, &githubService),
					resource.TestCheckResourceAttr(githubResourceName, "repository_url", "https://github.com/gitlabhq/terraform-provider-gitlab"),
					resource.TestCheckResourceAttr(githubResourceName, "static_context", "true"),
				),
			},
			// Update the github service
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabServiceGithubUpdateConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists(githubResourceName, &githubService),
					resource.TestCheckResourceAttr(githubResourceName, "repository_url", "https://github.com/terraform-providers/terraform-provider-github"),
					resource.TestCheckResourceAttr(githubResourceName, "static_context", "false"),
				),
			},
			// Update the github service to get back to previous settings
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabServiceGithubConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabServiceGithubExists(githubResourceName, &githubService),
					resource.TestCheckResourceAttr(githubResourceName, "repository_url", "https://github.com/gitlabhq/terraform-provider-gitlab"),
					resource.TestCheckResourceAttr(githubResourceName, "static_context", "true"),
				),
			},
		},
	})
}

// lintignore: AT002 // TODO: Resolve this tfproviderlint issue
func TestAccGitlabServiceGithub_import(t *testing.T) {
	githubResourceName := "gitlab_service_github.github"
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabServiceGithubDestroy,
		Steps: []resource.TestStep{
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabServiceGithubConfig(rInt),
			},
			{
				SkipFunc:          isRunningInCE,
				ResourceName:      githubResourceName,
				ImportStateIdFunc: getGithubProjectID(githubResourceName),
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
		conn := testAccProvider.Meta().(*gitlab.Client)

		githubService, _, err := conn.Services.GetGithubService(project)
		if err != nil {
			return fmt.Errorf("Github service does not exist in project %s: %v", project, err)
		}
		*service = *githubService

		return nil
	}
}

func testAccCheckGitlabServiceGithubDestroy(s *terraform.State) error {
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

func testAccGitlabServiceGithubConfig(rInt int) string {
	return fmt.Sprintf(`
# Requires billing with Silver or above
resource "gitlab_project" "foo" {
	name         = "foo-%d"
	description  = "Terraform acceptance tests"
}

resource "gitlab_service_github" "github" {
	project        = "${gitlab_project.foo.id}"
	token          = "test"
  repository_url = "https://github.com/gitlabhq/terraform-provider-gitlab"
}
`, rInt)
}

func testAccGitlabServiceGithubUpdateConfig(rInt int) string {
	return fmt.Sprintf(`
# Requires billing with Silver or above
resource "gitlab_project" "foo" {
	name         = "foo-%d"
	description  = "Terraform acceptance tests"
}

resource "gitlab_service_github" "github" {
	project        = "${gitlab_project.foo.id}"
	token          = "test"
	repository_url = "https://github.com/terraform-providers/terraform-provider-github"
	static_context = false
}
`, rInt)
}
