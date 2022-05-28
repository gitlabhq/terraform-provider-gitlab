package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabClusterAgentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project": {
			Description: "ID or full path of the project maintained by the authenticated user.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"name": {
			Description: "The Name of the agent.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"agent_id": {
			Description: "The ID of the agent.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"created_at": {
			Description: "The ISO8601 datetime when the agent was created.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"created_by_user_id": {
			Description: "The ID of the user who created the agent.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
	}
}

func gitlabClusterAgentToStateMap(project string, clusterAgent *gitlab.Agent) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["name"] = clusterAgent.Name
	stateMap["agent_id"] = clusterAgent.ID
	stateMap["created_at"] = clusterAgent.CreatedAt.Format(time.RFC3339)
	stateMap["created_by_user_id"] = clusterAgent.CreatedByUserID
	return stateMap
}
