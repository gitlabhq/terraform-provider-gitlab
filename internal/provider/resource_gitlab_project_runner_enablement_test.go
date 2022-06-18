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
	gitlab "github.com/xanzy/go-gitlab"
)

func TestAccGitlabProjectRunnerEnablement_basic(t *testing.T) {
	testGroup := testAccCreateGroups(t, 1)[0]
	projectA := testAccCreateProjectWithNamespace(t, testGroup.ID)
	projectB := testAccCreateProjectWithNamespace(t, testGroup.ID)
	projectC := testAccCreateProjectWithNamespace(t, testGroup.ID)

	name := fmt.Sprintf("TestAcc Runner %s", acctest.RandString(10))

	opts := gitlab.RegisterNewRunnerOptions{
		Token:       &projectA.RunnersToken,
		Description: gitlab.String(name),
	}

	// Create runner in project A
	runner, _, _ := testGitlabClient.Runners.RegisterNewRunner(&opts)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabProjectRunnerEnablementDestroy(projectB.ID, runner.ID),
		Steps: []resource.TestStep{
			// Enable it in projectB
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_runner_enablement" "foo" {
					project = %d
					runner_id = %d
				}`, projectB.ID, runner.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "project", fmt.Sprint(projectB.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "runner_id", fmt.Sprint(runner.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "id", fmt.Sprintf("%d:%d", projectB.ID, runner.ID)),
					// The runner is enabled in project B
					testAccCheckGitlabProjectRunnerEnablementCreate(projectB.ID, runner.ID),
				),
			},
			// Verify foo resource with an import.
			{
				ResourceName:      "gitlab_project_runner_enablement.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Enable it in projectC
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_runner_enablement" "foo" {
					project = %d
					runner_id = %d
				}`, projectC.ID, runner.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "project", fmt.Sprint(projectC.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "runner_id", fmt.Sprint(runner.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "id", fmt.Sprintf("%d:%d", projectC.ID, runner.ID)),
					testAccCheckGitlabProjectRunnerEnablementCreate(projectC.ID, runner.ID),
					// The runner is no longer enabled on B
					testAccCheckGitlabProjectRunnerEnablementDestroy(projectB.ID, runner.ID),
				),
			},
			// Verify foo resource with an import.
			{
				ResourceName:      "gitlab_project_runner_enablement.foo",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Enable it in both projects
			{
				Config: fmt.Sprintf(`
				resource "gitlab_project_runner_enablement" "foo" {
					project = %d
					runner_id = %d
				}

				resource "gitlab_project_runner_enablement" "bar" {
					project = %d
					runner_id = %d
				}`, projectC.ID, runner.ID, projectB.ID, runner.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "project", fmt.Sprint(projectC.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "runner_id", fmt.Sprint(runner.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.foo", "id", fmt.Sprintf("%d:%d", projectC.ID, runner.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.bar", "project", fmt.Sprint(projectB.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.bar", "runner_id", fmt.Sprint(runner.ID)),
					resource.TestCheckResourceAttr("gitlab_project_runner_enablement.bar", "id", fmt.Sprintf("%d:%d", projectB.ID, runner.ID)),
					testAccCheckGitlabProjectRunnerEnablementCreate(projectB.ID, runner.ID),
					testAccCheckGitlabProjectRunnerEnablementCreate(projectC.ID, runner.ID),
				),
			},
			// Verify bar resource with an import.
			{
				ResourceName:      "gitlab_project_runner_enablement.bar",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckGitlabProjectRunnerEnablementCreate(pid int, rid int) resource.TestCheckFunc {
	return func(_ *terraform.State) error {
		runnerdetails, _, err := testGitlabClient.Runners.GetRunnerDetails(rid)
		if err != nil {
			return err
		}

		for _, p := range runnerdetails.Projects {
			if p.ID == pid {
				// The runner is enabled in the project - no error
				return nil
			}
		}

		return errors.New("Runner is not enabled in the project")
	}
}

func testAccCheckGitlabProjectRunnerEnablementDestroy(pid int, rid int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		testCreate := testAccCheckGitlabProjectRunnerEnablementCreate(pid, rid)
		err := testCreate(s)
		if err.Error() != "Runner is not enabled in the project" {
			return err
		}
		return nil
	}
}
