package provider

import (
	"sort"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

func gitlabProjectIssueBoardSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project": {
			Description: "The ID or full path of the project maintained by the authenticated user.",
			Type:        schema.TypeString,
			ForceNew:    true,
			Required:    true,
		},
		"name": {
			Description: "The name of the board.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"assignee_id": {
			Description: "The assignee the board should be scoped to. Requires a GitLab EE license.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"milestone_id": {
			Description: "The milestone the board should be scoped to. Requires a GitLab EE license.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"labels": {
			Description: "The list of label names which the board should be scoped to. Requires a GitLab EE license.",
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
		"weight": {
			Description:      "The weight range from 0 to 9, to which the board should be scoped to. Requires a GitLab EE license.",
			Type:             schema.TypeInt,
			Optional:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntBetween(0, 9)),
		},
		"lists": {
			Description: "The list of issue board lists",
			Type:        schema.TypeList,
			Optional:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"id": {
						Description: "The ID of the list",
						Type:        schema.TypeInt,
						Computed:    true,
					},
					"label_id": {
						Description: "The ID of the label the list should be scoped to. Requires a GitLab EE license.",
						Type:        schema.TypeInt,
						Optional:    true,
						// NOTE(TF): not supported by the SDK yet, see https://github.com/hashicorp/terraform-plugin-sdk/issues/71
						//           Anyways, GitLab will complain about this, so no big deal ...
						// ConflictsWith: []string{"lists.assignee_id", "lists.milestone_id", "lists.iteration_id"},
					},
					"assignee_id": {
						Description: "The ID of the assignee the list should be scoped to. Requires a GitLab EE license.",
						Type:        schema.TypeInt,
						Optional:    true,
						// NOTE(TF): not supported by the SDK yet, see https://github.com/hashicorp/terraform-plugin-sdk/issues/71
						//           Anyways, GitLab will complain about this, so no big deal ...
						// ConflictsWith: []string{"lists.label_id", "lists.milestone_id", "lists.iteration_id"},
					},
					"milestone_id": {
						Description: "The ID of the milestone the list should be scoped to. Requires a GitLab EE license.",
						Type:        schema.TypeInt,
						Optional:    true,
						// NOTE(TF): not supported by the SDK yet, see https://github.com/hashicorp/terraform-plugin-sdk/issues/71
						//           Anyways, GitLab will complain about this, so no big deal ...
						// ConflictsWith: []string{"lists.label_id", "lists.assignee_id", "lists.iteration_id"},
					},
					"iteration_id": {
						Description: "The ID of the iteration the list should be scoped to. Requires a GitLab EE license.",
						Type:        schema.TypeInt,
						Optional:    true,
						// NOTE(TF): not supported by the SDK yet, see https://github.com/hashicorp/terraform-plugin-sdk/issues/71
						//           Anyways, GitLab will complain about this, so no big deal ...
						// ConflictsWith: []string{"lists.label_id", "lists.assignee_id", "lists.milestone_id"},
					},
					"position": {
						Description: "The position of the list within the board. The position for the list is based on the its position in the `lists` array.",
						Type:        schema.TypeInt,
						Computed:    true,
					},
				},
			},
		},
	}
}

func gitlabProjectIssueBoardToStateMap(project string, issueBoard *gitlab.IssueBoard) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["name"] = issueBoard.Name
	if issueBoard.Milestone != nil {
		stateMap["milestone_id"] = issueBoard.Milestone.ID
	} else {
		stateMap["milestone_id"] = nil
	}
	if issueBoard.Assignee != nil {
		stateMap["assignee_id"] = issueBoard.Assignee.ID
	} else {
		stateMap["assignee_id"] = nil
	}
	stateMap["weight"] = issueBoard.Weight
	stateMap["labels"] = extractLabelNames(issueBoard.Labels)
	stateMap["lists"] = flattenProjectIssueBoardLists(issueBoard.Lists)
	return stateMap
}

func flattenProjectIssueBoardLists(lists []*gitlab.BoardList) (values []map[string]interface{}) {
	// GitLab returns the lists in arbitrary order, so we need to sort them by position first
	sort.Slice(lists, func(i, j int) bool {
		return lists[i].Position < lists[j].Position
	})
	for _, list := range lists {
		v := map[string]interface{}{
			"id":       list.ID,
			"position": list.Position,
		}

		if list.Label != nil {
			v["label_id"] = list.Label.ID
		} else {
			v["label_id"] = nil
		}
		if list.Assignee != nil {
			v["assignee_id"] = list.Assignee.ID
		} else {
			v["assignee_id"] = nil
		}
		if list.Milestone != nil {
			v["milestone_id"] = list.Milestone.ID
		} else {
			v["milestone_id"] = nil
		}
		if list.Iteration != nil {
			v["iteration_id"] = list.Iteration.ID
		} else {
			v["iteration_id"] = nil
		}

		values = append(values, v)
	}
	return values
}

func extractLabelNames(labels []*gitlab.LabelDetails) []string {
	var labelNames []string
	for _, label := range labels {
		labelNames = append(labelNames, label.Name)
	}
	return labelNames
}
