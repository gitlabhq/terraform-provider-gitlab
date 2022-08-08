package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/xanzy/go-gitlab"
)

func gitlabProjectHookSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"project": {
			Description: "The name or id of the project to add the hook to.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"project_id": {
			Description: "The id of the project for the hook.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"hook_id": {
			Description: "The id of the project hook.",
			Type:        schema.TypeInt,
			Computed:    true,
		},
		"url": {
			Description: "The url of the hook to invoke.",
			Type:        schema.TypeString,
			Required:    true,
		},
		"token": {
			Description: "A token to present when invoking the hook. The token is not available for imported resources.",
			Type:        schema.TypeString,
			Optional:    true,
			Sensitive:   true,
		},
		"push_events": {
			Description: "Invoke the hook for push events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
		"push_events_branch_filter": {
			Description: "Invoke the hook for push events on matching branches only.",
			Type:        schema.TypeString,
			Optional:    true,
		},
		"issues_events": {
			Description: "Invoke the hook for issues events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"confidential_issues_events": {
			Description: "Invoke the hook for confidential issues events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"merge_requests_events": {
			Description: "Invoke the hook for merge requests.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"tag_push_events": {
			Description: "Invoke the hook for tag push events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"note_events": {
			Description: "Invoke the hook for notes events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"confidential_note_events": {
			Description: "Invoke the hook for confidential notes events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"job_events": {
			Description: "Invoke the hook for job events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"pipeline_events": {
			Description: "Invoke the hook for pipeline events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"wiki_page_events": {
			Description: "Invoke the hook for wiki page events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"deployment_events": {
			Description: "Invoke the hook for deployment events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"releases_events": {
			Description: "Invoke the hook for releases events.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     false,
		},
		"enable_ssl_verification": {
			Description: "Enable ssl verification when invoking the hook.",
			Type:        schema.TypeBool,
			Optional:    true,
			Default:     true,
		},
	}
}

func gitlabProjectHookToStateMap(project string, hook *gitlab.ProjectHook) map[string]interface{} {
	stateMap := make(map[string]interface{})
	stateMap["project"] = project
	stateMap["project_id"] = hook.ProjectID
	stateMap["hook_id"] = hook.ID
	stateMap["url"] = hook.URL
	stateMap["push_events"] = hook.PushEvents
	stateMap["push_events_branch_filter"] = hook.PushEventsBranchFilter
	stateMap["issues_events"] = hook.IssuesEvents
	stateMap["confidential_issues_events"] = hook.ConfidentialIssuesEvents
	stateMap["merge_requests_events"] = hook.MergeRequestsEvents
	stateMap["tag_push_events"] = hook.TagPushEvents
	stateMap["note_events"] = hook.NoteEvents
	stateMap["confidential_note_events"] = hook.ConfidentialNoteEvents
	stateMap["job_events"] = hook.JobEvents
	stateMap["pipeline_events"] = hook.PipelineEvents
	stateMap["wiki_page_events"] = hook.WikiPageEvents
	stateMap["deployment_events"] = hook.DeploymentEvents
	stateMap["releases_events"] = hook.ReleasesEvents
	stateMap["enable_ssl_verification"] = hook.EnableSSLVerification
	return stateMap
}
