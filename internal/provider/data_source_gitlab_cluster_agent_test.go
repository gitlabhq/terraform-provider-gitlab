//go:build acceptance
// +build acceptance

package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccDataSourceGitlabClusterAgent_basic(t *testing.T) {
	testAccRequiresAtLeast(t, "14.10")

	testProject := testAccCreateProject(t)
	testAgent := testAccCreateClusterAgents(t, testProject.ID, 1)[0]

	resource.ParallelTest(t, resource.TestCase{
		ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
					data "gitlab_cluster_agent" "this" {
						project           = "%d"
						agent_id          = %d
					}
					`, testProject.ID, testAgent.ID,
				),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "name", testAgent.Name),
					resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "created_at", testAgent.CreatedAt.Format(time.RFC3339)),
					resource.TestCheckResourceAttr("data.gitlab_cluster_agent.this", "created_by_user_id", fmt.Sprintf("%d", testAgent.CreatedByUserID)),
				),
			},
		},
	})
}
