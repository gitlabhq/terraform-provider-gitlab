package provider

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func gitlabReleaseLinkGetSchema() map[string]*schema.Schema {
	validLinkTypes := []string{"other", "runbook", "image", "package"}

	return map[string]*schema.Schema{
		"project": {
			Description: "The ID or [URL-encoded path of the project](https://docs.gitlab.com/ee/api/index.html#namespaced-path-encoding).",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"tag_name": {
			Description: "The tag associated with the Release.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"name": {
			Description: "The name of the link. Link names must be unique within the release.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"url": {
			Description: "The URL of the link. Link URLs must be unique within the release.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"filepath": {
			Description: "Relative path for a [Direct Asset link](https://docs.gitlab.com/ee/user/project/releases/index.html#permanent-links-to-release-assets).",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"link_type": {
			Description:      fmt.Sprintf("The type of the link. Valid values are %s. Defaults to %s.", renderValueListForDocs(validLinkTypes), validLinkTypes[0]),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          validLinkTypes[0],
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validLinkTypes, false)),
		},
		"link_id": {
			Description: "The ID of the link.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"direct_asset_url": {
			Description: "Full path for a [Direct Asset link](https://docs.gitlab.com/ee/user/project/releases/index.html#permanent-links-to-release-assets).",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"external": {
			Description: "External or internal link.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}

func gitlabReleaseLinkToStateMap(project string, tagName string, releaseLink *gitlab.ReleaseLink) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["tag_name"] = tagName
	stateMap["name"] = releaseLink.Name
	stateMap["url"] = releaseLink.URL
	directAssetLinkArray := strings.SplitN(releaseLink.DirectAssetURL, "downloads", 2)
	if len(directAssetLinkArray) > 1 {
		stateMap["filepath"] = directAssetLinkArray[1]
	} else {
		stateMap["filepath"] = ""
	}
	stateMap["link_type"] = releaseLink.LinkType
	stateMap["link_id"] = releaseLink.ID
	stateMap["direct_asset_url"] = releaseLink.DirectAssetURL
	stateMap["external"] = releaseLink.External

	return stateMap
}
