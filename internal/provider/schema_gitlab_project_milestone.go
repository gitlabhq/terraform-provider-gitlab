package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func gitlabProjectMilestoneGetSchema() map[string]*schema.Schema {
	validMilestoneStates := []string{"active", "closed"}

	return map[string]*schema.Schema{
		"iid": {
			Description: "The ID of the project’s milestone.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"milestone_id": {
			Description: "The ID of the project’s milestone in Gitlab DB.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"project_id": {
			Description: "The ID or URL-encoded path of the project owned by the authenticated user.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"title": {
			Description: "The title of a milestone.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"description": {
			Description: "The description of the milestone.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"due_date": {
			Description:      "The due date of the milestone.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: isISO6801Date,
		},
		"start_date": {
			Description:      "The start date of the milestone.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: isISO6801Date,
		},
		"updated_at": {
			Description: "The last update time of the milestone.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
			//       we use it here to avoid having to parse the string ourselves.
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
		},
		"created_at": {
			Description: "The time of creation of the milestone.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
			// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
			//       we use it here to avoid having to parse the string ourselves.
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
		},
		// NOTE: not part of `CREATE`, but part of `UPDATE` with the `state_event` field.
		"state": {
			Description:      fmt.Sprintf("The state of the milestone. Valid values are: %s.", renderValueListForDocs(validMilestoneStates)),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "active",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validMilestoneStates, false)),
		},
		"web_url": {
			Description: "The web URL of the milestone.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"expired": {
			Description: "Bool, true if milestore expired.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
	}
}

func gitlabProjectMilestoneToStateMap(milestone *gitlab.Milestone) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["iid"] = milestone.IID
	stateMap["milestone_id"] = milestone.ID
	stateMap["project_id"] = fmt.Sprintf("%d", milestone.ProjectID)
	stateMap["title"] = milestone.Title
	stateMap["description"] = milestone.Description
	if milestone.DueDate != nil {
		stateMap["due_date"] = milestone.DueDate.String()
	} else {
		stateMap["due_date"] = nil
	}
	if milestone.StartDate != nil {
		stateMap["start_date"] = milestone.StartDate.String()
	} else {
		stateMap["start_date"] = nil
	}
	stateMap["updated_at"] = milestone.UpdatedAt.Format(time.RFC3339)
	stateMap["created_at"] = milestone.CreatedAt.Format(time.RFC3339)
	stateMap["state"] = milestone.State
	stateMap["web_url"] = milestone.WebURL
	if milestone.Expired != nil {
		stateMap["expired"] = milestone.Expired
	} else {
		stateMap["expired"] = false
	}

	return stateMap
}