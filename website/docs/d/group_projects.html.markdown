---
layout: "gitlab"
page_title: "GitLab: gitlab_group_projects"
sidebar_current: "docs-gitlab-data-source-group_projects"
description: |-
  View information about all projects inside a group
---

# gitlab\_group\_projects

Provides a list of projects in a specific group

## Example Usage

```hcl
data "gitlab_group_projects" "example" {
    group_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of the group you want to view the projects for

* `order_by` - (Optional) Return projects ordered by id, name, path, created_at, updated_at, or last_activity_at fields. Default is created_at.

* `search` - (Optional) Return list of projects matching the search criteria

## Attributes Reference

The following attributes are exported:

* `projects` - A list containing the projects matching the supplied arguments

Projects have the following fields:

* `id` - The ID of the project

* `name` - The name of the project
