//go:build acceptance
// +build acceptance

package provider

import (
	"errors"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectProtectedEnvironment_basic(t *testing.T) {
	testAccCheckEE(t)

	// Set up project environment.
	project := testAccCreateProject(t)
	environment := testAccCreateProjectEnvironment(t, project.ID, &gitlab.CreateEnvironmentOptions{
		Name: gitlab.String(acctest.RandomWithPrefix("test-protected-environment")),
	})

	// Set up project user.
	user := testAccCreateUsers(t, 1)[0]
	testAccAddProjectMembers(t, project.ID, []*gitlab.User{user})

	// Set up group access.
	group := testAccCreateGroups(t, 1)[0]
	if _, err := testGitlabClient.Projects.ShareProjectWithGroup(project.ID, &gitlab.ShareWithGroupOptions{
		GroupID:     &group.ID,
		GroupAccess: gitlab.AccessLevel(gitlab.MaintainerPermissions),
	}); err != nil {
		t.Fatalf("unable to share project %d with group %d", project.ID, group.ID)
	}

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectProtectedEnvironmentDestroy(project.ID, environment.Name),
		Steps: []resource.TestStep{
			// Create a basic protected environment.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_protected_environment" "this" {
					project     = %d
					environment = %q
					deploy_access_levels {
						access_level = "developer"
					}
				}`, project.ID, environment.Name),
				// Check computed attributes.
				Check: resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.0.access_level_description"),
			},
			// Verify upstream attributes with an import.
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the protected environment.
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_protected_environment" "this" {
					project     = %d
					environment = %q
					required_approval_count = 1
					deploy_access_levels {
						access_level = "maintainer"
					}
					deploy_access_levels {
						user_id = %d
					}
					deploy_access_levels {
						group_id = %d
					}
				}`, project.ID, environment.Name, user.ID, group.ID),
				// Check computed attributes.
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.0.access_level_description"),
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.1.access_level_description"),
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.2.access_level_description"),
					// access_level is computed when not specified.
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.1.access_level"),
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "deploy_access_levels.2.access_level"),
					// required_approval_count is set.
					resource.TestCheckResourceAttrSet("gitlab_project_protected_environment.this", "required_approval_count"),
				),
			},
			// Verify upstream attributes with an import.
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectProtectedEnvironmentDestroy(projectID int, environmentName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		_, _, err := testGitlabClient.ProtectedEnvironments.GetProtectedEnvironment(projectID, environmentName)
		if err == nil {
			return errors.New("environment is still protected")
		}
		if !is404(err) {
			return fmt.Errorf("unable to get protected environment: %w", err)
		}
		return nil
	}
}
