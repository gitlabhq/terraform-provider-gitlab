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

func TestAccGitlabClusterAgentToken_basic(t *testing.T) {
	testAccRequiresAtLeast(t, "15.0")

	testProject := testAccCreateProject(t)
	testAgent := testAccCreateClusterAgents(t, testProject.ID, 1)[0]
	var sutClusterAgentToken gitlab.AgentToken

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		CheckDestroy:      testAccCheckGitlabClusterAgentTokenDestroy,
		Steps: []resource.TestStep{
			// Verify creation with minimal required attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_cluster_agent_token" "this" {
						project  = "%d"
						agent_id = "%d"
						name     = "agent-1-token"
					}
				`, testProject.ID, testAgent.ID),
				Check: resource.ComposeTestCheckFunc(
					// Verify that computed attributes have been correctly set to state
					testAccResourceGitlabClusterAgentTokenGet("gitlab_cluster_agent_token.this", &sutClusterAgentToken),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "token_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgentToken.ID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "created_at", func(value string) error {
						expectedValue := sutClusterAgentToken.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil

					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "created_by_user_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgentToken.CreatedByUserID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrSet("gitlab_cluster_agent_token.this", "token"),
				),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_cluster_agent_token.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
			// Verify update with all attributes
			{
				Config: fmt.Sprintf(`
					resource "gitlab_cluster_agent_token" "this" {
						project     = "%d"
						agent_id    = "%d"
						name        = "agent-1-token"
						description = "agent-1-description"
					}
				`, testProject.ID, testAgent.ID),
				Check: resource.ComposeTestCheckFunc(
					// Verify that computed attributes have been correctly set to state
					testAccResourceGitlabClusterAgentTokenGet("gitlab_cluster_agent_token.this", &sutClusterAgentToken),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "token_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgentToken.ID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "created_at", func(value string) error {
						expectedValue := sutClusterAgentToken.CreatedAt.Format(time.RFC3339)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil

					}),
					resource.TestCheckResourceAttrWith("gitlab_cluster_agent_token.this", "created_by_user_id", func(value string) error {
						expectedValue := fmt.Sprintf("%d", sutClusterAgentToken.CreatedByUserID)
						if value != expectedValue {
							return fmt.Errorf("should be equal to %s", expectedValue)
						}
						return nil
					}),
					resource.TestCheckResourceAttrSet("gitlab_cluster_agent_token.this", "token"),
				),
			},
			// Verify Import
			{
				ResourceName:            "gitlab_cluster_agent_token.this",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"token"},
			},
		},
	})
}

func testAccResourceGitlabClusterAgentTokenGet(resourceName string, clusterAgentToken *gitlab.AgentToken) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("resource '%q' not found", resourceName)
		}

		project, agentID, tokenID, err := resourceGitlabClusterAgentTokenParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject, _, err := testGitlabClient.ClusterAgents.GetAgentToken(project, agentID, tokenID)
		if err != nil {
			return err
		}

		*clusterAgentToken = *subject
		return nil
	}
}

func testAccCheckGitlabClusterAgentTokenDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gitlab_cluster_agent_token" {
			continue
		}

		project, agentID, tokenID, err := resourceGitlabClusterAgentTokenParseID(rs.Primary.ID)
		if err != nil {
			return err
		}

		subject, _, err := testGitlabClient.ClusterAgents.GetAgentToken(project, agentID, tokenID)
		if err == nil && subject != nil && subject.Status != "revoked" {
			return fmt.Errorf("gitlab_cluster_agent_token resource '%s' not yet revoked", rs.Primary.ID)
		}

		if err != nil && !is404(err) {
			return err
		}

		return nil
	}
	return nil
}
