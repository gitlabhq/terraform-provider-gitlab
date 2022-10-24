package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabUserKeysSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Description: "The ID or username of the group.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"title": {
			Description: "The title of the key",
			Type:        schema.TypeString,
			Required:    false,
		},
		"key": {
			Description: "The SSH key",
			Type:        schema.TypeString,
			Required:    false,
			ForceNew:    true,
		},
		"created_at": {
			Description: "The creation date for this key",
			Type:        schema.TypeString,
			Required:    false,
		},
	}
}

func gitlabUserKeyToStateMap(key *gitlab.SSHKey) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["id"] = key.ID
	stateMap["title"] = key.Title
	stateMap["key"] = key.Key
	stateMap["created_at"] = key.CreatedAt.Format(time.RFC3339)
	return stateMap
}
