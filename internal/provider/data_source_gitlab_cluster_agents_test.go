//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabClusterAgents_basic(t *testing.T) {
	testAccRequiresAtLeast(t, "14.10")

	testProject := testAccCreateProject(t)
	testClusterAgents := testAccCreateClusterAgents(t, testProject.ID, 25)

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_cluster_agents" "this" {
						project = "%d"
					}
				`, testProject.ID),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_cluster_agents.this", "cluster_agents.#", fmt.Sprintf("%d", len(testClusterAgents))),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.0.name"),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.0.created_at"),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.0.created_by_user_id"),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.1.name"),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.1.created_at"),
					resource.TestCheckResourceAttrSet("data.gitlab_cluster_agents.this", "cluster_agents.1.created_by_user_id"),
				),
			},
		},
	})
}
