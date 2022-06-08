//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/xanzy/go-gitlab"
)

func TestAccGitlabClusterAgent_basic(t *testing.T) {
	testAccRequiresAtLeast(t, "14.10")

	testProject := testAccCreateProject(t)
	var sutClusterAgent gitlab.Agent

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabClusterAgentDestroy,
		Steps: []resource.TestStep{
			// Verify creation
			{
				Config: fmt.Sprintf(`
					resource "gitlab_cluster_agent" "this" {
						project = "%d"
						name    = "agent-1"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					// Verify that computed attributes have been correctly set to state
					testAccResourceGitlabClusterAgentGet("gitlab_cluster_agent.this", &sutClusterAgent),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "agent_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgent.ID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "created_at", func(value string) error {
						expectedValue := sutClusterAgent.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil

					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "created_by_user_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgent.CreatedByUserID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_cluster_agent.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Verify update / re-creation
			{
				Config: fmt.Sprintf(`
					resource "gitlab_cluster_agent" "this" {
						project = "%d"
						name    = "agent-2"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					// Verify that computed attributes have been correctly set to state
					testAccResourceGitlabClusterAgentGet("gitlab_cluster_agent.this", &sutClusterAgent),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "agent_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgent.ID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "created_at", func(value string) error {
						expectedValue := sutClusterAgent.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil

					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent.this", "created_by_user_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgent.CreatedByUserID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
				),
			},
			// Verify Import
			{
				ResourceName:      "gitlab_cluster_agent.this",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccResourceGitlabClusterAgentGet(resourceName string, clusterAgent *gitlab.Agent) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource '%q' not found", resourceName)
		}

		project, agentID, err := resourceGitlabClusterAgentParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject, _, err := testGitlabClient.ClusterAgents.GetAgent(project, agentID)
		if err != nil {
			return err
		}

		*clusterAgent = *subject
		return nil
	}
}

func testAccCheckGitlabClusterAgentDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_cluster_agent" {
			continue
		}

		project, agentID, err := resourceGitlabClusterAgentParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject, _, err := testGitlabClient.ClusterAgents.GetAgent(project, agentID)
		if err == nil && subject != nil {
			return fmt.Errorf("gitlab_cluster_agent resource '%s' still exists", rs.Primary.ID)
		}

		if err != nil && !is404(err) {
			return err
		}

		return nil
	}
	return nil
}
