package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func gitlabProjectTagGetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description: "The name of a tag.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"message": {
			Description: "The message of the annotated tag.",
			Type:        schema.TypeString,
			Optional:    true,
			ForceNew:    true,
		},
		"protected": {
			Description: "Bool, true if tag has tag protection.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"target": {
			Description: "The unique id assigned to the commit by Gitlab.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"release": {
			Description: "The release associated with the tag.",
			Type:        schema.TypeSet,
			Computed:    true,
			Set:         schema.HashResource(releaseNoteSchema),
			Elem:        releaseNoteSchema,
		},
		"commit": {
			Description: "The commit associated with the tag.",
			Type:        schema.TypeSet,
			Computed:    true,
			Set:         schema.HashResource(commitSchema),
			Elem:        commitSchema,
		},
	}
}

var releaseNoteSchema = &schema.Resource{
	Schema: map[string]*schema.Schema{
		"tag_name": {
			Description: "The name of the tag.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"description": {
			Description: "The description of release.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	},
}
