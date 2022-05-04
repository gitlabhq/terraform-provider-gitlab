package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_cluster_agent", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_cluster_agent`" + ` data source allows to retrieve details about a GitLab Agent for Kubernetes.

-> Requires at least GitLab 14.10

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/cluster_agents.html)`,

		ReadContext: dataSourceGitlabClusterAgentRead,
		Schema:      datasourceSchemaFromResourceSchema(gitlabClusterAgentSchema(), []string{"project", "agent_id"}, nil),
	}
})

func dataSourceGitlabClusterAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	agentID := d.Get("agent_id").(int)

	clusterAgent, _, err := client.ClusterAgents.GetAgent(project, agentID, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s:%d", project, agentID))
	stateMap := gitlabClusterAgentToStateMap(project, clusterAgent)
	if err := setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}
