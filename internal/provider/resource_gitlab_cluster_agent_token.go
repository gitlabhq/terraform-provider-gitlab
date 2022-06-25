package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_cluster_agent_token", func() *schema.Resource {
	return &schema.Resource{
		Description: `The ` + "`" + `gitlab_cluster_agent_token` + "`" + ` resource allows to manage the lifecycle of a token for a GitLab Agent for Kubernetes.

-> Requires at least maintainer permissions on the project.

-> Requires at least GitLab 15.0

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/cluster_agents.html#create-an-agent-token)`,

		CreateContext: resourceGitlabClusterAgentTokenCreate,
		ReadContext:   resourceGitlabClusterAgentTokenRead,
		DeleteContext: resourceGitlabClusterAgentTokenDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: gitlabClusterAgentTokenSchema(),
	}
})

func resourceGitlabClusterAgentTokenCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)

	project := d.Get("project").(string)
	agentID := d.Get("agent_id").(int)
	options := gitlab.CreateAgentTokenOptions{
		Name: gitlab.String(d.Get("name").(string)),
	}

	if v, ok := d.GetOk("description"); ok {
		options.Description = gitlab.String(v.(string))
	}

	log.Printf("[DEBUG] create token for GitLab Agent for Kubernetes %d in project %s with name '%v'", agentID, project, options.Name)
	clusterAgentToken, _, err := client.ClusterAgents.CreateAgentToken(project, agentID, &options, gitlab.WithContext(ctx))
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resourceGitlabClusterAgentTokenBuildID(project, agentID, clusterAgentToken.ID))
	// NOTE: the token is only returned with the direct response from the create API.
	d.Set("token", clusterAgentToken.Token)
	return resourceGitlabClusterAgentTokenRead(ctx, d, meta)
}

func resourceGitlabClusterAgentTokenRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, agentID, tokenID, err := resourceGitlabClusterAgentTokenParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] read token for GitLab Agent for Kubernetes %d in project %s with id %d", agentID, project, tokenID)
	clusterAgentToken, _, err := client.ClusterAgents.GetAgentToken(project, agentID, tokenID, gitlab.WithContext(ctx))
	if err != nil {
		if is404(err) {
			log.Printf("[DEBUG] read token for GitLab Agent for Kubernetes %d in project %s with id %d not found, removing from state", agentID, project, tokenID)
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	stateMap := gitlabClusterAgentTokenToStateMap(project, clusterAgentToken)
	if err = setStateMapInResourceData(stateMap, d); err != nil {
		return diag.FromErr(err)
	}
	return nil
}

func resourceGitlabClusterAgentTokenDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*gitlab.Client)
	project, agentID, tokenID, err := resourceGitlabClusterAgentTokenParseID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[DEBUG] delete token for GitLab Agent for Kubernetes %d in project %s with id %d", agentID, project, tokenID)
	if _, err := client.ClusterAgents.RevokeAgentToken(project, agentID, tokenID, gitlab.WithContext(ctx)); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceGitlabClusterAgentTokenBuildID(project string, agentID int, tokenID int) string {
	return fmt.Sprintf("%s:%d:%d", project, agentID, tokenID)
}

func resourceGitlabClusterAgentTokenParseID(id string) (string, int, int, error) {
	parts := strings.Split(id, ":")
	if len(parts) != 3 {
		return "", 0, 0, fmt.Errorf("invalid cluster agent token id %q, expected format '{project}:{agent_id}:{token_id}", id)
	}
	project, rawAgentID, rawTokenID := parts[0], parts[1], parts[2]
	agentID, err := strconv.Atoi(rawAgentID)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid cluster agent token id %q with 'agent_id' %q, expected integer", id, rawAgentID)
	}
	tokenID, err := strconv.Atoi(rawTokenID)
	if err != nil {
		return "", 0, 0, fmt.Errorf("invalid cluster agent token id %q with 'token_id' %q, expected integer", id, rawTokenID)
	}

	return project, agentID, tokenID, nil
}
