---
layout: "gitlab"
page_title: "GitLab: gitlab_service_jira"
sidebar_current: "docs-gitlab-resource-service_jira"
description: |-
  Manage Jira integration that allows to receive event notifications in Jira
---

# gitlab\_service\_jira

This resource allows you to manage Jira integration.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_jira" "jira" {
  project  = "${gitlab_project.awesome_project.id}"
  url      = "https://jira.example.com"
  username = "user"
  password = "mypass"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) ID of the project you want to activate integration on.

* `url` - (Required) The URL to the JIRA project which is being linked to this GitLab project. For example, https://jira.example.com.

* `username` - (Required) The username of the user created to be used with GitLab/JIRA.

* `password` - (Required) The password of the user created to be used with GitLab/JIRA.

* `project_key` - (Required) The short identifier for your JIRA project, all uppercase, e.g., PROJ.

* `jira_issue_transition_id` - (Optional) The ID of a transition that moves issues to a closed state. You can find this number under the JIRA workflow administration (Administration > Issues > Workflows) by selecting View under Operations of the desired workflow of your project. By default, this ID is set to 2.

## Importing Jira service

 You can import a service_jira state using `terraform import <resource> <project_id>`:

```bash
$ terraform import gitlab_service_jira.jira 1
```