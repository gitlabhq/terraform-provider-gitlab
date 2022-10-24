package provider

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"group_id": {
			Description: "",
			Type:        schema.TypeInt,
			Required:    true,
		},
		"name": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"path": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"description": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"visibility": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"share_with_group_lock": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"require_two_factor_authentication": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"two_factor_grace_period": {
			Description: "",
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
		},
		"project_creation_level": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"auto_devops_enabled": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"subgroup_creation_level": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"emails_disabled": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"mentions_disabled": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"lfs_enabled": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"default_branch_protection": {
			Description: "",
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
		},
		"avatar_url": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"web_url": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"request_access_enabled": {
			Description: "",
			Type:        schema.TypeBool,
			Computed:    true,
			Optional:    true,
		},
		"full_name": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"full_path": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"file_template_project_id": {
			Description: "",
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
		},
		"parent_id": {
			Description: "",
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
		},
		"created_at": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"ip_restriction_ranges": {
			Description: "",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"statistics": {
			Description: "",
			Type:        schema.TypeMap,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Computed: true,
			Optional: true,
		},
	}
}

func gitlabGroupToStateMap(group *gitlab.Group) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["group_id"] = group.ID
	stateMap["auto_devops_enabled"] = group.AutoDevopsEnabled
	stateMap["avatar_url"] = group.AvatarURL
	stateMap["created_at"] = group.CreatedAt.Format(time.RFC3339)
	stateMap["default_branch_protection"] = group.DefaultBranchProtection
	stateMap["description"] = group.Description
	stateMap["emails_disabled"] = group.EmailsDisabled
	stateMap["file_template_project_id"] = group.FileTemplateProjectID
	stateMap["full_name"] = group.FullName
	stateMap["full_path"] = group.FullPath
	stateMap["group_id"] = group.ID
	stateMap["lfs_enabled"] = group.LFSEnabled
	stateMap["mentions_disabled"] = group.MentionsDisabled
	stateMap["name"] = group.Name
	stateMap["parent_id"] = group.ParentID
	stateMap["path"] = group.Path
	stateMap["project_creation_level"] = group.ProjectCreationLevel
	stateMap["request_access_enabled"] = group.RequestAccessEnabled
	stateMap["require_two_factor_authentication"] = group.RequireTwoFactorAuth
	stateMap["share_with_group_lock"] = group.ShareWithGroupLock
	stateMap["statistics"] = group.Statistics
	stateMap["subgroup_creation_level"] = group.SubGroupCreationLevel
	stateMap["two_factor_grace_period"] = group.TwoFactorGracePeriod
	stateMap["visibility"] = group.Visibility
	stateMap["web_url"] = group.WebURL
	return stateMap
}
