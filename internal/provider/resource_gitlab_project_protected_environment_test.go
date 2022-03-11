package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectProtectedEnvironment_basic(t *testing.T) {

	var pt gitlab.ProtectedEnvironment
	rInt := acctest.RandInt()
	testProject := testAccCreateProject(t)
	testEnvironment := testAccCreateProjectEnvironment(t, testProject.ID, &gitlab.CreateEnvironmentOptions{
		Name: gitlab.String(fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt)),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectProtectedEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Create a Protected Environment with default options
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "maintainer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.MaintainerPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment to get back to initial settings
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectProtectedEnvironment_user(t *testing.T) {

	var pt gitlab.ProtectedEnvironment
	rInt := acctest.RandInt()
	testProject := testAccCreateProject(t)
	testUser := testAccCreateUsers(t, 1)[0]
	testAccAddProjectMembers(t, testProject.ID, []*gitlab.User{testUser})
	testEnvironment := testAccCreateProjectEnvironment(t, testProject.ID, &gitlab.CreateEnvironmentOptions{
		Name: gitlab.String(fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt)),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectProtectedEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Create a Protected Environment with default options
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment with user level access
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentUpdateUserConfig(testEnvironment.Name, testProject.ID, testUser.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:   fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						UserID: testUser.ID,
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment to get back to initial settings
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabProjectProtectedEnvironment_group(t *testing.T) {

	var pt gitlab.ProtectedEnvironment
	rInt := acctest.RandInt()
	testProject := testAccCreateProject(t)
	group := testAccCreateGroups(t, 1)[0]
	testEnvironment := testAccCreateProjectEnvironment(t, testProject.ID, &gitlab.CreateEnvironmentOptions{
		Name: gitlab.String(fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt)),
	})

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectProtectedEnvironmentDestroy,
		Steps: []resource.TestStep{
			// Create a Protected Environment with default options
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment with group level access
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentUpdateGroupConfig(testEnvironment.Name, testProject.ID, group.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:    fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						GroupID: group.ID,
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the Protected Environment to get back to initial settings
			{
				SkipFunc: isRunningInCE,
				Config:   testAccGitlabProjectProtectedEnvironmentBasicConfig(testEnvironment.Name, testProject.ID, "developer"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabProjectProtectedEnvironmentExists("gitlab_project_protected_environment.this", &pt),
					testAccCheckGitlabProjectProtectedEnvironmentAttributes(&pt, &testAccGitlabProjectProtectedEnvironmentExpectedAttributes{
						Name:              fmt.Sprintf("ProjectProtectedEnvironment-%d", rInt),
						CreateAccessLevel: accessLevelValueToName[gitlab.DeveloperPermissions],
					}),
				),
			},
			// Verify import
			{
				ResourceName:      "gitlab_project_protected_environment.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectProtectedEnvironmentExists(n string, pt *gitlab.ProtectedEnvironment) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		project, environment, err := parseTwoPartID(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("Error in Splitting Project ID and Environment Name")
		}

		pes, _, err := testGitlabClient.ProtectedEnvironments.ListProtectedEnvironments(project, nil)
		if err != nil {
			return err
		}

		for _, gotpe := range pes {
			if gotpe.Name == environment {
				*pt = *gotpe
				return nil
			}
		}

		return fmt.Errorf("Protected Environment does not exist")
	}
}

type testAccGitlabProjectProtectedEnvironmentExpectedAttributes struct {
	Name              string
	CreateAccessLevel string
	GroupID           int
	UserID            int
}

func testAccCheckGitlabProjectProtectedEnvironmentAttributes(pt *gitlab.ProtectedEnvironment, want *testAccGitlabProjectProtectedEnvironmentExpectedAttributes) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if pt.Name != want.Name {
			return fmt.Errorf("got name %q; want %q", pt.Name, want.Name)
		}

		if pt.DeployAccessLevels[0].AccessLevel != accessLevelNameToValue[want.CreateAccessLevel] {
			return fmt.Errorf("got create access levels %q; want %q", pt.DeployAccessLevels[0].AccessLevel, accessLevelNameToValue[want.CreateAccessLevel])
		}

		if pt.DeployAccessLevels[0].GroupID != want.GroupID {
			return fmt.Errorf("got group ID %q; want %q", pt.DeployAccessLevels[0].GroupID, want.GroupID)
		}

		if pt.DeployAccessLevels[0].UserID != want.UserID {
			return fmt.Errorf("got user ID %q; want %q", pt.DeployAccessLevels[0].UserID, want.UserID)
		}

		return nil
	}
}

func testAccCheckGitlabProjectProtectedEnvironmentDestroy(s *terraform.State) error {
	var project string
	var protectedEnvironment string
	var err error
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "gitlab_project" {
			project = rs.Primary.ID
		} else if rs.Type == "gitlab_project_protected_environment" {
			project, protectedEnvironment, err = parseTwoPartID(rs.Primary.ID)
			if err != nil {
				return fmt.Errorf("[ERROR] cannot get project and Protected Environment from input: %v", rs.Primary.ID)
			}
		}
	}

	pt, _, err := testGitlabClient.ProtectedEnvironments.GetProtectedEnvironment(project, protectedEnvironment)
	if err == nil {
		if pt != nil {
			return fmt.Errorf("project Protected Environment %s still exists", protectedEnvironment)
		}
	}

	if is404(err) {
		return err
	}

	return nil
}

func testAccGitlabProjectProtectedEnvironmentBasicConfig(environmentName string, projectID int, access_level string) string {
	return fmt.Sprintf(`
resource "gitlab_project_protected_environment" "this" {
  project     = %[2]d
  environment = "%[1]s"

  deploy_access_levels {
		access_level = "%[3]s"
  }
}
	`, environmentName, projectID, access_level)
}

func testAccGitlabProjectProtectedEnvironmentUpdateGroupConfig(environmentName string, projectID int, groupID int) string {
	return fmt.Sprintf(`
resource "gitlab_project_protected_environment" "this" {
  project     = %[2]d
  environment = "%[1]s"

  deploy_access_levels {
    group_id = %[3]d
  }
}
	`, environmentName, projectID, groupID)
}

func testAccGitlabProjectProtectedEnvironmentUpdateUserConfig(environmentName string, projectID int, userID int) string {
	return fmt.Sprintf(`
resource "gitlab_project_protected_environment" "this" {
  project     = %[2]d
  environment = "%[1]s"

  deploy_access_levels {
    user_id = %[3]d
  }
}
	`, environmentName, projectID, userID)
}
