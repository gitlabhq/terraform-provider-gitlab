package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_cluster_agent", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_cluster_agent` + "`" + ` resource allows to manage the lifecycle of a GitLab Agent for Kubernetes.

-> Note that this resource only registers the agent, but doesn't configure it.
   The configuration needs to be manually added as described in
   [the docs](https://docs.gitlab.com/ee/user/clusters/agent/install/index.html#create-an-agent-configuration-file).
   However, a ` + "`gitlab_repository_file`" + ` resource may be used to achieve that.

-> Requires at least maintainer permissions on the project.

-> Requires at least GitLab 14.10

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/cluster_agents.html)`,

		CreateContext: resourceGitlabClusterAgentCreate,
		ReadContext:   resourceGitlabClusterAgentRead,
		DeleteContext: resourceGitlabClusterAgentDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: gitlabClusterAgentSchema(),
	}
})

func resourceGitlabClusterAgentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	options := gitlab.RegisterAgentOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	log.Printf("[DEBUG] create GitLab Agent for Kubernetes in project %s with name '%v'", project, options.Name)
	clusterAgent, _, err := client.ClusterAgents.RegisterAgent(project, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGitlabClusterAgentBuildID(project, clusterAgent.ID))
	return resourceGitlabClusterAgentRead(ctx, d, meta)
}

func resourceGitlabClusterAgentRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, agentID, err := resourceGitlabClusterAgentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read GitLab Agent for Kubernetes in project %s with id %d", project, agentID)
	clusterAgent, _, err := client.ClusterAgents.GetAgent(project, agentID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] read GitLab Agent for Kubernetes in project %s with id %d not found, removing from state", project, agentID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabClusterAgentToStateMap(project, clusterAgent)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabClusterAgentDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, agentID, err := resourceGitlabClusterAgentParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete GitLab Agent for Kubernetes in project %s with id %d", project, agentID)
	if _, err := client.ClusterAgents.DeleteAgent(project, agentID, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabClusterAgentBuildID(project string, agentID int) string {
	return fmt.Sprintf("%s:%d", project, agentID)
}

func resourceGitlabClusterAgentParseID(id string) (string, int, error) {
	project, rawAgentID, err := parseTwoPartID(id)
	if err != nil {
		return "", 0, err
	}

	agentID, err := strconv.Atoi(rawAgentID)
	if err != nil {
		return "", 0, err
	}

	return project, agentID, nil
}
