package provider

import (
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func gitlabUserSSHKeySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"user_id": {
			Description: "The ID or username of the user.",
			Type:        schema.TypeInt,
			ForceNew:    true,
			Required:    true,
		},
		"title": {
			Description: "The title of the ssh key.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"key": {
			Description: "The ssh key. The SSH key `comment` (trailing part) is optional and ignored for diffing, because GitLab overrides it with the username and GitLab hostname.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
			DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
				// NOTE: the ssh keys consist of three parts: `type`, `data`, `comment`, whereas the `comment` is optional
				//       and suppressed in the diffing. It's overridden by GitLab with the username and GitLab hostname.
				newParts := strings.Fields(new)
				oldParts := strings.Fields(old)
				if len(newParts) < 2 || len(oldParts) < 2 {
					// NOTE: at least one of the keys doesn't have the required two parts, thus we just compare them
					return new == old
				}

				// NOTE: both keys have the required two parts, thus we compare the parts separately, ignoring the rest
				return newParts[0] == oldParts[0] && newParts[1] == oldParts[1]
			},
		},
		"expires_at": {
			Description: "The expiration date of the SSH key in ISO 8601 format (YYYY-MM-DDTHH:MM:SSZ)",
			Type:        schema.TypeString,
			Default:     nil,
			ForceNew:    true,
			Optional:    true,
			// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
			//       we use it here to avoid having to parse the string ourselves.
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
		},
		"key_id": {
			Description: "The ID of the ssh key.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"created_at": {
			Description: "The time when this key was created in GitLab.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func gitlabUserKeyToStateMap(key *gitlab.SSHKey) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["title"] = key.Title
	stateMap["key"] = key.Key
	stateMap["key_id"] = key.ID
	if key.ExpiresAt != nil {
		stateMap["expires_at"] = key.ExpiresAt.Format(time.RFC3339)
	}
	if key.CreatedAt != nil {
		stateMap["created_at"] = key.CreatedAt.Format(time.RFC3339)
	}
	return stateMap
}
