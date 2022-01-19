# gitlab\_service\_microsoft\_teams

This resource allows you to manage Microsoft Teams integration.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name             = "awesome_project"
  description      = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_microsoft_teams" "teams" {
  project     = gitlab_project.awesome_project.id
  webhook     = "https://testurl.com/?token=XYZ"
  push_events = true 
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) ID of the project you want to activate integration on.
* `webhook` - (Required) The Microsoft Teams webhook. For example, https://outlook.office.com/webhook/...
* `notify_only_broken_pipelines` - (Optional) Send notifications for broken pipelines
* `branches_to_be_notified` - (Optional) Branches to send notifications for. Valid options are “all”, “default”, “protected”, and “default_and_protected”. The default value is “default”
* `push_events` - (Optional) Enable notifications for push events
* `issues_events` - (Optional) Enable notifications for issue events
* `confidential_issues_events` - (Optional) Enable notifications for confidential issue events
* `merge_requests_events` - (Optional) Enable notifications for merge request events
* `tag_push_events` - (Optional) Enable notifications for tag push events
* `note_events` - (Optional) Enable notifications for note events
* `confidential_note_events` - (Optional) Enable notifications for confidential note events
* `pipeline_events` - (Optional) Enable notifications for pipeline events
* `wiki_page_events` - (Optional) Enable notifications for wiki page events

## Importing Microsoft Teams service

 You can import a service_microsoft_teams state using `terraform import <resource> <project_id>`:

```bash
$ terraform import gitlab_service_microsoft_teams.teams 1
```
