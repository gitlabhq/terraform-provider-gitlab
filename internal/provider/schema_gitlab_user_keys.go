package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabUserSSHKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"title": {
			Description: "The title of the key",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"key": {
			Description: "The SSH key",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"created_at": {
			Description: "The creation date for this key",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func gitlabUserKeyToStateMap(key *gitlab.SSHKey) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["title"] = key.Title
	stateMap["key"] = key.Key
	stateMap["created_at"] = key.CreatedAt.Format(time.RFC3339)
	return stateMap
}
