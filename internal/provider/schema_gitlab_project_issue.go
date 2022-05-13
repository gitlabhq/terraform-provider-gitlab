package provider

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/xanzy/go-gitlab"
)

var validIssueTypes = []string{"issue", "incident", "test_case"}

func gitlabProjectIssueGetSchema() map[string]*schema.Schema {
	validIssueStates := []string{"opened", "closed"}

	return map[string]*schema.Schema{
		"project": {
			Description: "The name or ID of the project.",
			Type:        schema.TypeString,
			Required:    true,
			ForceNew:    true,
		},
		"iid": {
			Description: "The internal ID of the project's issue.",
			Type:        schema.TypeInt,
			Optional:    true,
			ForceNew:    true,
			Computed:    true,
		},
		"title": {
			Description: "The title of the issue.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"issue_id": {
			Description: "The instance-wide ID of the issue.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		// NOTE: not supported yet in go-gitlab (v0.55.0)
		// "assignee_id": {
		// 	Description: "The ID of the user to assign the issue to. Only appears on GitLab Free.",
		// 	Type:        schema.TypeInt,
		// 	Optional:    true,
		// },
		"assignee_ids": {
			Description: "The IDs of the users to assign the issue to.",
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeInt},
			Optional:    true,
		},
		"confidential": {
			Description: "Set an issue to be confidential.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"created_at": {
			Description: "When the issue was created. Date time string, ISO 8601 formatted, for example 2016-03-11T03:45:40Z. Requires administrator or project/group owner rights.",
			Type:        schema.TypeString,
			// NOTE: cannot be updated
			ForceNew: true,
			Optional: true,
			Computed: true,
			// NOTE: since RFC3339 is pretty much a subset of ISO8601 and actually expected by GitLab,
			//       we use it here to avoid having to parse the string ourselves.
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsRFC3339Time),
		},
		"description": {
			Description: "The description of an issue. Limited to 1,048,576 characters.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"discussion_to_resolve": {
			Description: "The ID of a discussion to resolve. This fills out the issue with a default description and mark the discussion as resolved. Use in combination with merge_request_to_resolve_discussions_of.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"due_date": {
			Description: `The due date. Date time string in the format YYYY-MM-DD, for example 2016-03-11.
**Note:** removing a due date is currently not supported, see https://github.com/xanzy/go-gitlab/issues/1384 for details.
				`,
			Type:             schema.TypeString,
			Optional:         true,
			ValidateDiagFunc: isISO6801Date,
		},
		"epic_id": {
			Description: "ID of the epic to add the issue to. Valid values are greater than or equal to 0.",
			Type:        schema.TypeInt,
			// NOTE: not yet supported to be set in go-gitlab.
			Computed: true,
			// ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		"issue_type": {
			Description:      fmt.Sprintf("The type of issue. Valid values are: %s.", renderValueListForDocs(validIssueTypes)),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "issue",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueTypes, false)),
		},
		"labels": {
			Description: "The labels of an issue.",
			Type:        schema.TypeSet,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Optional:    true,
		},
		"merge_request_to_resolve_discussions_of": {
			Description: "The IID of a merge request in which to resolve all issues. This fills out the issue with a default description and mark all discussions as resolved. When passing a description or title, these values take precedence over the default values.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"milestone_id": {
			Description: "The global ID of a milestone to assign issue. To find the milestone_id associated with a milestone, view an issue with the milestone assigned and use the API to retrieve the issue's details.",
			Type:        schema.TypeInt,
			Optional:    true,
		},
		"weight": {
			Description:      "The weight of the issue. Valid values are greater than or equal to 0.",
			Type:             schema.TypeInt,
			Optional:         true,
			Computed:         true,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IntAtLeast(0)),
		},
		"external_id": {
			Description: "The external ID of the issue.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		// NOTE: not part of `CREATE`, but part of `UPDATE` with the `state_event` field.
		"state": {
			Description:      fmt.Sprintf("The state of the issue. Valid values are: %s.", renderValueListForDocs(validIssueStates)),
			Type:             schema.TypeString,
			Optional:         true,
			Default:          "opened",
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice(validIssueStates, false)),
		},
		// NOTE: to keep things simple, users of this resource should use the `gitlab_user` data source to
		//       get more information about the author if desired.
		"author_id": {
			Description: "The ID of the author of the issue. Use `gitlab_user` data source to get more information about the user.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"updated_at": {
			Description: "When the issue was updated. Date time string, ISO 8601 formatted, for example 2016-03-11T03:45:40Z.",
			Type:        schema.TypeString,
			Computed:    true,
			Optional:    true,
		},
		"closed_at": {
			Description: "When the issue was closed. Date time string, ISO 8601 formatted, for example 2016-03-11T03:45:40Z.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		// NOTE: to keep things simple, users of this resource should use the `gitlab_user` data source to
		//       get more information about the closer if desired.
		"closed_by_user_id": {
			Description: "The ID of the user that closed the issue. Use `gitlab_user` data source to get more information about the user.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"moved_to_id": {
			Description: "The ID of the issue that was moved to.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"upvotes": {
			Description: "The number of upvotes the issue has received.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"downvotes": {
			Description: "The number of downvotes the issue has received.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"web_url": {
			Description: "The web URL of the issue.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"references": {
			Description: "The references of the issue.",
			Type:        schema.TypeMap,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
		},
		// NOTE(TF): these are from the `time_stats` field.
		//           Clarify what to do with nested types.
		"time_estimate": {
			Description: "The time estimate of the issue.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"total_time_spent": {
			Description: "The total time spent of the issue.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"human_time_estimate": {
			Description: "The human-readable time estimate of the issue.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		"human_total_time_spent": {
			Description: "The human-readable total time spent of the issue.",
			Type:        schema.TypeString,
			Computed:    true,
		},
		// NOTE(TF): end `time_stats`
		// NOTE: not part of `CREATE`, but part of `UPDATE`
		"discussion_locked": {
			Description: "Whether the issue is locked for discussions or not.",
			Type:        schema.TypeBool,
			Optional:    true,
		},
		"subscribed": {
			Description: "Whether the authenticated user is subscribed to the issue or not.",
			Type:        schema.TypeBool,
			Computed:    true,
		},
		"user_notes_count": {
			Description: "The number of user notes on the issue.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"links": {
			Description: "The links of the issue.",
			Type:        schema.TypeMap,
			Elem:        &schema.Schema{Type: schema.TypeString},
			Computed:    true,
		},
		"issue_link_id": {
			Description: "The ID of the issue link.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"merge_requests_count": {
			Description: "The number of merge requests associated with the issue.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"epic_issue_id": {
			Description: "The ID of the epic issue.",
			Type:        schema.TypeInt,
			Computed:    true,
			Optional:    true,
		},
		"task_completion_status": {
			Description: "The task completion status. It's always a one element list.",
			Type:        schema.TypeList,
			Computed:    true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"count": {
						Description: "The number of tasks.",
						Type:        schema.TypeInt,
						Computed:    true,
						Optional:    true,
					},
					"completed_count": {
						Description: "The number of tasks that are completed.",
						Type:        schema.TypeInt,
						Computed:    true,
						Optional:    true,
					},
				},
			},
		},
	}
}

func gitlabProjectIssueToStateMap(project string, issue *gitlab.Issue) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["iid"] = issue.IID
	stateMap["title"] = issue.Title
	stateMap["issue_id"] = issue.ID
	stateMap["assignee_ids"] = flattenAssigneeIds(issue.Assignees)
	stateMap["confidential"] = issue.Confidential
	stateMap["created_at"] = issue.CreatedAt.Format(time.RFC3339)
	stateMap["description"] = issue.Description
	if issue.DueDate != nil {
		stateMap["due_date"] = issue.DueDate.String()
	} else {
		stateMap["due_date"] = nil
	}
	stateMap["issue_type"] = issue.IssueType
	stateMap["labels"] = issue.Labels
	if issue.Milestone != nil {
		stateMap["milestone_id"] = issue.Milestone.ID
	} else {
		stateMap["milestone_id"] = nil
	}
	stateMap["weight"] = issue.Weight
	stateMap["external_id"] = issue.ExternalID
	stateMap["state"] = issue.State
	if issue.Author != nil {
		stateMap["author_id"] = issue.Author.ID
	} else {
		stateMap["author_id"] = nil
	}
	if issue.UpdatedAt != nil {
		stateMap["updated_at"] = issue.UpdatedAt.Format(time.RFC3339)
	} else {
		stateMap["updated_at"] = nil
	}
	if issue.ClosedAt != nil {
		stateMap["closed_at"] = issue.ClosedAt.Format(time.RFC3339)
	} else {
		stateMap["closed_at"] = nil
	}
	if issue.ClosedBy != nil {
		stateMap["closed_by_user_id"] = issue.ClosedBy.ID
	} else {
		stateMap["closed_by_user_id"] = nil
	}
	stateMap["moved_to_id"] = issue.MovedToID
	stateMap["upvotes"] = issue.Upvotes
	stateMap["downvotes"] = issue.Downvotes
	stateMap["web_url"] = issue.WebURL
	if issue.References != nil {
		stateMap["references"] = flattenIssueReferences(issue.References)
	} else {
		stateMap["references"] = nil
	}
	if issue.TimeStats != nil {
		stateMap["time_estimate"] = issue.TimeStats.TimeEstimate
		stateMap["total_time_spent"] = issue.TimeStats.TotalTimeSpent
		stateMap["human_time_estimate"] = issue.TimeStats.HumanTimeEstimate
		stateMap["human_total_time_spent"] = issue.TimeStats.HumanTotalTimeSpent
	} else {
		stateMap["time_estimate"] = nil
		stateMap["total_time_spent"] = nil
		stateMap["human_time_estimate"] = nil
		stateMap["human_total_time_spent"] = nil
	}
	stateMap["discussion_locked"] = issue.DiscussionLocked
	stateMap["subscribed"] = issue.Subscribed
	stateMap["user_notes_count"] = issue.UserNotesCount
	if issue.Links != nil {
		stateMap["links"] = flattenIssueLinks(issue.Links)
	} else {
		stateMap["links"] = nil
	}
	stateMap["issue_link_id"] = issue.IssueLinkID
	stateMap["merge_requests_count"] = issue.MergeRequestCount
	stateMap["epic_issue_id"] = issue.EpicIssueID
	if issue.Epic != nil {
		stateMap["epic_id"] = issue.Epic.ID
	} else {
		stateMap["epic_id"] = nil
	}
	if issue.TaskCompletionStatus != nil {
		stateMap["task_completion_status"] = flattenIssueTaskCompletionStatus(issue.TaskCompletionStatus)
	} else {
		stateMap["task_completion_status"] = nil
	}

	return stateMap
}

func flattenAssigneeIds(assignees []*gitlab.IssueAssignee) (result []int) {
	if assignees == nil {
		return
	}

	for _, assignee := range assignees {
		result = append(result, assignee.ID)
	}
	return result
}

func flattenIssueReferences(references *gitlab.IssueReferences) (result map[string]string) {
	if references == nil {
		return
	}

	result = map[string]string{
		"short":    references.Short,
		"relative": references.Relative,
		"full":     references.Full,
	}
	return result
}

func flattenIssueLinks(issueLinks *gitlab.IssueLinks) (result map[string]string) {
	if issueLinks == nil {
		return
	}

	result = map[string]string{
		"self":        issueLinks.Self,
		"notes":       issueLinks.Notes,
		"award_emoji": issueLinks.AwardEmoji,
		"project":     issueLinks.Project,
	}
	return result
}

func flattenIssueTaskCompletionStatus(taskCompletionStatus *gitlab.TasksCompletionStatus) (result []map[string]interface{}) {
	if taskCompletionStatus == nil {
		return
	}

	result = []map[string]interface{}{
		{
			"count":           taskCompletionStatus.Count,
			"completed_count": taskCompletionStatus.CompletedCount,
		},
	}
	return result
}
