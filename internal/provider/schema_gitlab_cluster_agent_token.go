package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabClusterAgentTokenSchema() map[string]*schema.Schema {
	tokenStatuses := []string{"active", "revoked"}

	return map[string]*schema.Schema{
		"project": {
			Description: "ID or full path of the project maintained by the authenticated user.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"agent_id": {
			Description: "The ID of the agent.",
			Type:        schema.TypeInt,
			ForceNew:    true,
			Required:    true,
		},
		"token_id": {
			Description: "The ID of the token.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"name": {
			Description: "The Name of the agent.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"description": {
			Description: "The Description for the agent.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Optional:    true,
		},
		"status": {
			Description: fmt.Sprintf("The status of the token. Valid values are %s.", renderValueListForDocs(tokenStatuses)),
			Type:        schema.TypeString,
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
		"last_used_at": {
			Description: "The ISO8601 datetime when the token was last used.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"token": {
			Description: "The secret token for the agent. The `token` is not available in imported resources.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func gitlabClusterAgentTokenToStateMap(project string, clusterAgentToken *gitlab.AgentToken) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["agent_id"] = clusterAgentToken.AgentID
	stateMap["token_id"] = clusterAgentToken.ID
	stateMap["name"] = clusterAgentToken.Name
	stateMap["description"] = clusterAgentToken.Description
	stateMap["status"] = clusterAgentToken.Status
	stateMap["created_at"] = clusterAgentToken.CreatedAt.Format(time.RFC3339)
	stateMap["created_by_user_id"] = clusterAgentToken.CreatedByUserID
	if clusterAgentToken.LastUsedAt != nil {
		stateMap["last_used_at"] = clusterAgentToken.LastUsedAt.Format(time.RFC3339)
	}
	return stateMap
}
