---
layout: "gitlab"
page_title: "GitLab: gitlab_service_slack"
sidebar_current: "docs-gitlab-resource-service_slack"
description: |-
  Manage Slack notifications integration that allows to receive event notifications in Slack
---

# gitlab\_service_slack

This resource allows you to manage Slack notifications integration.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_slack" "slack" {
  project                    = "${gitlab_project.awesome_project.id}"
  webhook                    = "https://webhook.com"
  username                   = "myuser"
  push_events                = true
  push_channel               = "push_chan"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) ID of the project you want to activate integration on.

* `webhook` - (Required) Webhook URL (ex.: https://hooks.slack.com/services/...)

* `username` - (Optional) Username to use.

* `notify_only_broken_pipelines` - (Optional) Send notifications for broken pipelines.

* `notify_only_default_branch` - (Optional) Send notifications only for the default branch.

* `push_events` - (Optional) Enable notifications for push events.

* `push_channel` - (Optional) The name of the channel to receive push events notifications.

* `issues_events` - (Optional) Enable notifications for issues events.

* `issue_channel` - (Optional) The name of the channel to receive issue events notifications.

* `confidential_issues_events` - (Optional) Enable notifications for confidential issues events.

* `confidential_issue_channel` - (Optional) The name of the channel to receive confidential issue events notifications.

* `merge_requests_events` - (Optional) Enable notifications for merge requests events.

* `merge_request_channel` - (Optional) The name of the channel to receive merge request events notifications.

* `tag_push_events` - (Optional) Enable notifications for tag push events.

* `tag_push_channel` - (Optional) The name of the channel to receive tag push events notifications.

* `note_events` - (Optional) Enable notifications for note events.

* `note_channel` - (Optional) The name of the channel to receive note events notifications.

* `confidential_note_events` - (Optional) Enable notifications for confidential note events.

* `pipeline_events` - (Optional) Enable notifications for pipeline events.

* `pipeline_channel` - (Optional) The name of the channel to receive pipeline events notifications.

* `wiki_page_events` - (Optional) Enable notifications for wiki page events.

* `wiki_page_channel` - (Optional) The name of the channel to receive wiki page events notifications.

* `deployment_events` - (Optional) Enable notifications for deployment events.

* `deployment_channel` - (Optional) The name of the channel to receive deployment events notifications.

## Importing Slack service

You can import a service_slack state using `terraform import <resource> <project_id>`:

```bash
$ terraform import gitlab_service_slack.slack 1
```
