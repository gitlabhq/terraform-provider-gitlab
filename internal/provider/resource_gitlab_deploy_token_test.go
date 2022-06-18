//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabDeployToken_basic(t *testing.T) {
	var projectDeployToken gitlab.DeployToken
	var groupDeployToken gitlab.DeployToken

	testProject := testAccCreateProject(t)
	testGroup := testAccCreateGroups(t, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabDeployTokenConfig(testProject.ID, testGroup.ID),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckGitlabDeployTokenExists("gitlab_deploy_token.project_token", &projectDeployToken),
					resource.TestCheckResourceAttrSet("gitlab_deploy_token.project_token", "token"),
					testAccCheckGitlabDeployTokenExists("gitlab_deploy_token.group_token", &groupDeployToken),
					resource.TestCheckResourceAttrSet("gitlab_deploy_token.group_token", "token"),
				),
			},
			// Verify import
			{
				ResourceName:            "gitlab_deploy_token.project_token",
				ImportStateIdFunc:       getDeployTokenImportID("gitlab_deploy_token.project_token"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			{
				ResourceName:            "gitlab_deploy_token.group_token",
				ImportStateIdFunc:       getDeployTokenImportID("gitlab_deploy_token.group_token"),
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}
func TestAccGitlabDeployToken_pagination(t *testing.T) {
	testGroup := testAccCreateGroups(t, 1)[0]
	testProject := testAccCreateProject(t)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployTokenDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccGitlabDeployTokenPaginationConfig(25, testGroup.ID, testProject.ID),
			},
			// In case pagination wouldn't properly work, we would get that the plan isn't empty,
			// because some of the deploy tokens wouldn't be in the first page and therefore
			// considered non-existing, ...
			{
				Config:   testAccGitlabDeployTokenPaginationConfig(25, testGroup.ID, testProject.ID),
				PlanOnly: true,
			},
		},
	})
}

func testAccCheckGitlabDeployTokenExists(n string, deployToken *gitlab.DeployToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not Found: %s", n)
		}

		deployTokenID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		projectName := rs.Primary.Attributes["project"]
		groupName := rs.Primary.Attributes["group"]

		var gotDeployTokens []*gitlab.DeployToken

		if projectName != "" {
			gotDeployTokens, _, err = testGitlabClient.DeployTokens.ListProjectDeployTokens(projectName, nil)
		} else if groupName != "" {
			gotDeployTokens, _, err = testGitlabClient.DeployTokens.ListGroupDeployTokens(groupName, nil)
		} else {
			return fmt.Errorf("No project or group ID is set")
		}

		if err != nil {
			return err
		}

		for _, token := range gotDeployTokens {
			if token.ID == deployTokenID {
				*deployToken = *token
				return nil
			}
		}

		return fmt.Errorf("Deploy Token doesn't exist")
	}
}

func getDeployTokenImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		deployTokenID := rs.Primary.ID
		if deployTokenID == "" {
			return "", fmt.Errorf("No deploy token ID is set")
		}
		projectID := rs.Primary.Attributes["project"]
		if projectID != "" {
			return fmt.Sprintf("project:%s:%s", projectID, deployTokenID), nil
		}
		groupID := rs.Primary.Attributes["group"]
		if groupID != "" {
			return fmt.Sprintf("group:%s:%s", groupID, deployTokenID), nil
		}

		return "", fmt.Errorf("No project or group ID is set")
	}
}

func testAccCheckGitlabDeployTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_deploy_token" {
			continue
		}

		deployTokenID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}

		project := rs.Primary.Attributes["project"]
		group := rs.Primary.Attributes["group"]

		var gotDeployTokens []*gitlab.DeployToken

		if project != "" {
			gotDeployTokens, _, err = testGitlabClient.DeployTokens.ListProjectDeployTokens(project, nil)
		} else if group != "" {
			gotDeployTokens, _, err = testGitlabClient.DeployTokens.ListGroupDeployTokens(group, nil)
		} else {
			return fmt.Errorf("somehow neither project nor group were set")
		}

		if err == nil {
			for _, token := range gotDeployTokens {
				if token.ID == deployTokenID {
					return fmt.Errorf("Deploy token still exists")
				}
			}
		}

		if !is404(err) {
			return err
		}
	}

	return nil
}

func testAccGitlabDeployTokenConfig(projectID int, groupID int) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_token" "project_token" {
  project  = "%d"
  name     = "project-deploy-token"
  username = "my-username"

  expires_at = "2021-03-14T07:20:50.000Z"

  scopes = [
	"read_registry",
	"read_repository",
	"read_package_registry",
	"write_registry",
	"write_package_registry",
  ]
}

resource "gitlab_deploy_token" "group_token" {
  group  = "%d"
  name     = "group-deploy-token"
  username = "my-username"

  expires_at = "2021-03-14T07:20:50.000Z"

  scopes = [
	"read_registry",
	"read_repository",
	"read_package_registry",
	"write_registry",
	"write_package_registry",
  ]
}
  `, projectID, groupID)
}

func testAccGitlabDeployTokenPaginationConfig(numberOfTokens int, groupID int, projectID int) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_token" "example_group" {
  group  = %d
  name   = "deploy-token-${count.index}"
  scopes = ["read_registry"]

  count = %d
}

resource "gitlab_deploy_token" "example_project" {
  project  = %d
  name   = "deploy-token-${count.index}"
  scopes = ["read_registry"]

  count = %d
}
  `, groupID, numberOfTokens, projectID, numberOfTokens)
}

type expiresAtSuppressFuncTest struct {
	description string
	old         string
	new         string
	expected    bool
}

func TestExpiresAtSuppressFunc(t *testing.T) {
	t.Parallel()

	testcases := []expiresAtSuppressFuncTest{
		{
			description: "same dates without millis",
			old:         "2025-03-14T00:00:00Z",
			new:         "2025-03-14T00:00:00Z",
			expected:    true,
		}, {
			description: "different date without millis",
			old:         "2025-03-14T00:00:00Z",
			new:         "2025-03-14T11:11:11Z",
			expected:    false,
		}, {
			description: "same date with and without millis",
			old:         "2025-03-14T00:00:00Z",
			new:         "2025-03-14T00:00:00.000Z",
			expected:    true,
		}, {
			description: "cannot parse new date",
			old:         "2025-03-14T00:00:00Z",
			new:         "invalid-date",
			expected:    false,
		},
	}

	for _, test := range testcases {
		t.Run(test.description, func(t *testing.T) {
			actual := expiresAtSuppressFunc("", test.old, test.new, nil)
			if actual != test.expected {
				t.Fatalf("FAIL\n\told: %s, new: %s\n\texpected: %t\n\tactual: %t",
					test.old, test.new, test.expected, actual)
			}
		})
	}
}
