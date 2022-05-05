package provider

import (
	"context"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

var _ = registerDataSource("gitlab_cluster_agents", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`gitlab_cluster_agents`" + ` data source allows details of GitLab Agents for Kubernetes in a project.

-> Requires at least GitLab 14.10

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/cluster_agents.html)`,

		ReadContext: dataSourceGitlabClusterAgentsRead,
		Schema: map[string]*schema.Schema{
			"project": {
				Description: "The ID or full path of the project owned by the authenticated user.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"cluster_agents": {
				Description: "List of the registered agents.",
				Type:        schema.TypeList,
				Computed:    true,
				Elem: &schema.Resource{
					Schema: datasourceSchemaFromResourceSchema(gitlabClusterAgentSchema(), nil, nil),
				},
			},
		},
	}
})

func dataSourceGitlabClusterAgentsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.ListAgentsOptions{
		PerPage: 20,
		Page:    1,
	}

	var clusterAgents []*gitlab.Agent
	for options.Page != 0 {
		paginatedClusterAgents, resp, err := client.ClusterAgents.ListAgents(project, &options, gitlab.WithContext(ctx))
		if err != nil {
			return diag.FromErr(err)
		}

		clusterAgents = append(clusterAgents, paginatedClusterAgents...)
		options.Page = resp.NextPage
	}

	log.Printf("[DEBUG] list GitLab Agents for Kubernetes in project %s", project)
	d.SetId(project)
	d.Set("project", project)
	if err := d.Set("cluster_agents", flattenClusterAgentsForState(clusterAgents)); err != nil {
		return diag.Errorf("Failed to set cluster agents to state: %v", err)
	}
	return nil
}

func flattenClusterAgentsForState(clusterAgents []*gitlab.Agent) (values []map[string]interface{}) {
	for _, clusterAgent := range clusterAgents {
		values = append(values, map[string]interface{}{
			"name":               clusterAgent.Name,
			"created_at":         clusterAgent.CreatedAt.Format(time.RFC3339),
			"created_by_user_id": clusterAgent.CreatedByUserID,
		})
	}
	return values
}
