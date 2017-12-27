---
layout: "gitlab"
page_title: "GitLab: gitlab_project_membership"
sidebar_current: "docs-gitlab-resource-project-membership"
description: |-
  Adds a user to a project as a member
---

# gitlab\_project_membership

This resource allows you to add a current user to an existing project with a set access level.

## Example Usage

```hcl
resource "gitlab_project_membership" "test" {
project_id = "Test Project"
user_id = "Test User"
access_level = "guest"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The id of the project.

* `user_id` - (Required) The id of the user.

* `access_level` - (Required) One of five levels of access to the project.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the membership by the GitLab server.

