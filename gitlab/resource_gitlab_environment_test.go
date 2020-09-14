package gitlab

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabEnvironment_basic(t *testing.T) {
	var environment gitlab.Environment
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckGitlabEnvironmentDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabEnvironmentConfig(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabEnvironmentExists("gitlab_environment.environment", &environment),
					testAccCheckGitlabEnvironmentAttributes(&environment, &testAccGitlabEnvironmentAttributes{
						Name:        "meow",
						ExternalURL: "https://google.com",
					}),
				),
			},
			{
				Config: testAccGitlabEnvironmentConfigNameOnly(rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabEnvironmentExists("gitlab_environment.environment", &environment),
					testAccCheckGitlabEnvironmentAttributes(&environment, &testAccGitlabEnvironmentAttributes{
						Name: "meow",
					}),
				),
			},
		},
	})
}

func testAccCheckGitlabEnvironmentExists(n string, environment *gitlab.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		environmentID := rs.Primary.ID
		project := rs.Primary.Attributes["project"]
		if project == "" {
			return fmt.Errorf("No project ID is set")
		}

		conn := testAccProvider.Meta().(*gitlab.Client)

		environments, _, err := conn.Environments.ListEnvironments(project, nil)
		if err != nil {
			return err
		}

		for _, gotEnvironment := range environments {
			resourceID := fmt.Sprintf("%s/%d", project, gotEnvironment.ID)
			if resourceID == environmentID {
				*environment = *gotEnvironment
				return nil
			}
		}
		return fmt.Errorf("Environment does not exist")
	}
}

type testAccGitlabEnvironmentAttributes struct {
	Name        string
	ExternalURL string
}

func testAccCheckGitlabEnvironmentAttributes(environment *gitlab.Environment, want *testAccGitlabEnvironmentAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if environment.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", environment.Name, want.Name)
		}

		if environment.ExternalURL != want.ExternalURL {
			return fmt.Errorf("got external url %q; want %q", environment.ExternalURL, want.ExternalURL)
		}

		return nil
	}
}

func testAccCheckGitlabEnvironmentDestroy(s *terraform.State) error {
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

func testAccGitlabEnvironmentConfig(rInt int) string {
	return fmt.Sprintf(`
	resource "gitlab_project" "foo" {
		name = "foo-%d"
		description = "Terraform acceptance tests"
	  
		# So that acceptance tests can be run in a gitlab organization
		# with no billing
		visibility_level = "public"
	  }
	  
	  resource "gitlab_environment" "environment" {
		  project = "${gitlab_project.foo.id}"
		  name = "meow"
		  external_url = "https://google.com"
	  }
		  `, rInt)
}

func testAccGitlabEnvironmentConfigNameOnly(rInt int) string {
	return fmt.Sprintf(`
	resource "gitlab_project" "foo" {
		name = "foo-%d"
		description = "Terraform acceptance tests"
	  
		# So that acceptance tests can be run in a gitlab organization
		# with no billing
		visibility_level = "public"
	  }
	  
	  resource "gitlab_environment" "environment" {
		  project = "${gitlab_project.foo.id}"
		  name = "meow"
	  }
		  `, rInt)
}
