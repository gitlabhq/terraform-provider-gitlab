package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectEnvironment_basic(t *testing.T) {
	rInt := acctest.RandInt()

	testProject := testAccCreateProject(t)

	var env1 gitlab.Environment = gitlab.Environment{
		Name: fmt.Sprintf("ProjectEnvironment-%d", rInt),
	}

	var env2 gitlab.Environment = gitlab.Environment{
		Name:        fmt.Sprintf("ProjectEnvironment-%d", rInt),
		ExternalURL: "https://example.com",
	}

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Create an Environment with default options
			{
				Config: testAccGitlabProjectEnvironmentConfig(testProject.ID, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectEnvironmentExists("gitlab_project_environment.this", &env1),
					testAccCheckGitlabProjectEnvironmentAttributes(&env1, &testAccGitlabProjectEnvironmentExpectedAttributes{
						Name: fmt.Sprintf("ProjectEnvironment-%d", rInt),
					}),
					testCheckResourceAttrLazy("gitlab_project_environment.this", "created_at", func() string { return env1.CreatedAt.Format(time.RFC3339) }),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Environment
			{
				Config: testAccGitlabProjectEnvironmentUpdateConfig(testProject.ID, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectEnvironmentExists("gitlab_project_environment.this", &env2),
					testAccCheckGitlabProjectEnvironmentAttributes(&env2, &testAccGitlabProjectEnvironmentExpectedAttributes{
						Name:        fmt.Sprintf("ProjectEnvironment-%d", rInt),
						ExternalURL: "https://example.com",
					}),
					testCheckResourceAttrLazy("gitlab_project_environment.this", "created_at", func() string { return env2.CreatedAt.Format(time.RFC3339) }),
					testCheckResourceAttrLazy("gitlab_project_environment.this", "updated_at", func() string { return env2.UpdatedAt.Format(time.RFC3339) }),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Environment to get back to initial settings
			{
				Config: testAccGitlabProjectEnvironmentConfig(testProject.ID, rInt),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectEnvironmentExists("gitlab_project_environment.this", &env1),
					testAccCheckGitlabProjectEnvironmentAttributes(&env1, &testAccGitlabProjectEnvironmentExpectedAttributes{
						Name: fmt.Sprintf("ProjectEnvironment-%d", rInt),
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectEnvironment_stop_before_destroy(t *testing.T) {
	rInt := acctest.RandInt()

	testProject := testAccCreateProject(t)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Create environment with required values only
			{
				Config:      testAccGitlabProjectEnvironmentConfig(testProject.ID, rInt),
				ExpectError: regexp.MustCompile("Environment must be in a stopped state before deletion"),
			},
		},
	})
}

func testAccCheckGitlabProjectEnvironmentExists(n string, env *gitlab.Environment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project, environment, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in Splitting Project ID and Environment Name")
		}

		environmentID, err := strconv.Atoi(environment)
		if err != nil {
			return fmt.Errorf("error converting environment ID to int: %v", err)
		}

		if _, _, err := testGitlabClient.Environments.GetEnvironment(project, environmentID); err != nil {
			return err
		}
		return nil
	}
}

type testAccGitlabProjectEnvironmentExpectedAttributes struct {
	Name        string
	ExternalURL string
	State       string
}

func testAccCheckGitlabProjectEnvironmentAttributes(env *gitlab.Environment, want *testAccGitlabProjectEnvironmentExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if env.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", env.Name, want.Name)
		}

		if env.ExternalURL != want.ExternalURL {
			return fmt.Errorf("got external URL %q; want %q", env.ExternalURL, want.ExternalURL)
		}

		if env.State != want.State {
			return fmt.Errorf("got State %q; want %q", env.State, want.State)
		}

		return nil
	}
}

func testAccCheckGitlabProjectEnvironmentDestroy(s *terraform.State) error {
	var project string
	var environmentIDString string
	var environmentIDInt int
	var err error
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_project_environment" {
			project, environmentIDString, err = parseTwoPartID(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] cannot get project and environmentID from input: %v", rs.Primary.ID)
			}

			environmentIDInt, err = strconv.Atoi(environmentIDString)
			if err != nil {
				return fmt.Errorf("[ERROR] cannot convert environment ID to int: %v", err)
			}
		}
	}

	env, _, err := testGitlabClient.Environments.GetEnvironment(project, environmentIDInt)
	if err == nil {
		if env != nil {
			return fmt.Errorf("[ERROR] project Environment %v still exists", environmentIDInt)
		}
	}

	if is404(err) {
		return err
	}

	return nil
}

func testAccGitlabProjectEnvironmentConfig(projectID int, rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_environment" "this" {
  project = %d
  name    = "ProjectEnvironment-%d"

	stop_before_destroy = true
}
`, projectID, rInt)
}

func testAccGitlabProjectEnvironmentUpdateConfig(projectID int, rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_environment" "this" {
  project      = %d
  name         = "ProjectEnvironment-%d"
  external_url = "https://example.com"
}
`, projectID, rInt)
}

func testAccGitlabProjectEnvironmentStopBeforeDestroyFalse(projectID int, rInt int) string {
	return fmt.Sprintf(`
resource "gitlab_project_environment" "this" {
  project = %d
  name    = "ProjectEnvironment-%d"

  stop_before_destroy = false
}
`, projectID, rInt)
}
