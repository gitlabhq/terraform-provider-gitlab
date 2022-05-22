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
		"project": {
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
			Description:      "The due date of the milestone. Date time string in the format YYYY-MM-DD, for example 2016-03-11.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: isISO6801Date,
		},
		"start_date": {
			Description:      "The start date of the milestone. Date time string in the format YYYY-MM-DD, for example 2016-03-11.",
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: isISO6801Date,
		},
		// NOTE: not part of `CREATE`, but part of `UPDATE` with the `state_event` field.
		"state": {
			Description:      fmt.Sprintf("The state of the milestone. Valid values are: %s.", renderValueListForDocs(validMilestoneStates)),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "active",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validMilestoneStates, false)),
		},
		"created_at": {
			Description: "The time of creation of the milestone. Date time string, ISO 8601 formatted, for example 2016-03-11T03:45:40Z.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"expired": {
			Description: "Bool, true if milestore expired.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"iid": {
			Description: "The ID of the project's milestone.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"milestone_id": {
			Description: "The instance-wide ID of the project’s milestone.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"project_id": {
			Description: "The project ID of milestone.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"updated_at": {
			Description: "The last update time of the milestone. Date time string, ISO 8601 formatted, for example 2016-03-11T03:45:40Z.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"web_url": {
			Description: "The web URL of the milestone.",
			Type:        schema.TypeString,
			Computed:    true,
		},
	}
}

func gitlabProjectMilestoneToStateMap(project string, milestone *gitlab.Milestone) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["iid"] = milestone.IID
	stateMap["milestone_id"] = milestone.ID
	stateMap["project"] = project
	stateMap["project_id"] = milestone.ProjectID
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
	if milestone.UpdatedAt != nil {
		stateMap["updated_at"] = milestone.UpdatedAt.Format(time.RFC3339)
	} else {
		stateMap["updated_at"] = nil
	}
	if milestone.CreatedAt != nil {
		stateMap["created_at"] = milestone.CreatedAt.Format(time.RFC3339)
	} else {
		stateMap["created_at"] = nil
	}
	stateMap["state"] = milestone.State
	stateMap["web_url"] = milestone.WebURL
	if milestone.Expired != nil {
		stateMap["expired"] = milestone.Expired
	} else {
		stateMap["expired"] = false
	}

	return stateMap
}
