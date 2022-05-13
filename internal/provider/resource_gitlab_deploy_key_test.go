//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccGitlabDeployKey_basic(t *testing.T) {
	testProject := testAccCreateProject(t)
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployKeyDestroy,
		Steps: []resource.TestStep{
			// Create a project and deployKey with default options
			{
				Config: testAccGitlabDeployKeyConfig(rInt, "", testProject.ID),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key.foo",
				ImportStateIdFunc: getDeployKeyImportID("gitlab_deploy_key.foo"),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the project deployKey to toggle all the values to their inverse
			{
				Config: testAccGitlabDeployKeyUpdateConfig(rInt, testProject.ID),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key.foo",
				ImportStateIdFunc: getDeployKeyImportID("gitlab_deploy_key.foo"),
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Update the project deployKey to toggle the options back
			{
				Config: testAccGitlabDeployKeyConfig(rInt, "", testProject.ID),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key.foo",
				ImportStateIdFunc: getDeployKeyImportID("gitlab_deploy_key.foo"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccGitlabDeployKey_suppressTrailingSpace(t *testing.T) {
	testProject := testAccCreateProject(t)
	rInt := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabDeployKeyDestroy,
		Steps: []resource.TestStep{
			// Create a project and deployKey with space as suffix
			{
				Config: testAccGitlabDeployKeyConfig(rInt, " ", testProject.ID),
			},
			// Verify import
			{
				ResourceName:      "gitlab_deploy_key.foo",
				ImportStateIdFunc: getDeployKeyImportID("gitlab_deploy_key.foo"),
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabDeployKeyDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		deployKeyID, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return fmt.Errorf("unable to convert GitLab deploy key string into an int: %w", err)
		}

		project := rs.Primary.Attributes["project"]

		gotDeployKey, _, err := testGitlabClient.DeployKeys.GetDeployKey(project, deployKeyID)
		if err == nil {
			if gotDeployKey != nil && fmt.Sprintf("%d", gotDeployKey.ID) == rs.Primary.ID {
				return fmt.Errorf("Deploy key still exists")
			}
		}
		if !is404(err) {
			return err
		}
		return nil
	}
	return nil
}

func getDeployKeyImportID(n string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return "", fmt.Errorf("Not Found: %s", n)
		}

		deployKeyID := rs.Primary.ID
		if deployKeyID == "" {
			return "", fmt.Errorf("No deploy key ID is set")
		}

		projectID := rs.Primary.Attributes["project"]
		if projectID == "" {
			return "", fmt.Errorf("No project ID is set")
		}

		return fmt.Sprintf("%s:%s", projectID, deployKeyID), nil
	}
}

func testAccGitlabDeployKeyConfig(rInt int, suffix string, projectId int) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_key" "foo" {
  project = %[3]d
  title = "deployKey-%[1]d"
  key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCj13ozEBZ0s4el4k6mYqoyIKKKMh9hHY0sAYqSPXs2zGuVFZss1P8TPuwmdXVjHR7TiRXwC49zDrkyWJgiufggYJ1VilOohcMOODwZEJz+E5q4GCfHuh90UEh0nl8B2R0Uoy0LPeg93uZzy0hlHApsxRf/XZJz/1ytkZvCtxdllxfImCVxJReMeRVEqFCTCvy3YuJn0bce7ulcTFRvtgWOpQsr6GDK8YkcCCv2eZthVlrEwy6DEpAKTRiRLGgUj4dPO0MmO4cE2qD4ualY01PhNORJ8Q++I+EtkGt/VALkecwFuBkl18/gy+yxNJHpKc/8WVVinDeFrd/HhiY9yU0d richardc@tamborine.example.1%[2]s"
}
  `, rInt, suffix, projectId)
}

func testAccGitlabDeployKeyUpdateConfig(rInt int, projectId int) string {
	return fmt.Sprintf(`
resource "gitlab_deploy_key" "foo" {
  project = %[2]d
  title = "modifiedDeployKey-%[1]d"
  key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC6pSke2kb7YBjo65xDKegbOQsAtnMupRcFxXji7L1iXivGwORq0qpC2xzbhez5jk1WgPckEaNv2/Bz0uEW6oSIXw1KT1VN2WzEUfQCbpNyZPtn4iV3nyl6VQW/Nd1SrxiFJtH1H4vu+eCo4McMXTjuBBD06fiJNrHaSw734LjQgqtXWJuVym9qS5MqraZB7wDwTQwSM6kslL7KTgmo3ONsTLdb2zZhv6CS+dcFKinQo7/ttTmeMuXGbPOVuNfT/bePVIN1MF1TislHa2L2dZdGeoynNJT4fVPjA2Xl6eHWh4ySbvnfPznASsjBhP0n/QKprYJ/5fQShdBYBcuQiIMd richardc@tamborine.example.2"
  can_push = true
}
  `, rInt, projectId)
}
