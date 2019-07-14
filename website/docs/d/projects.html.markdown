---
layout: "gitlab"
page_title: "GitLab: gitlab_projects"
sidebar_current: "docs-gitlab-data-source-projects"
description: |-
  View information about multiple projects
---

# gitlab\_projects

Provides details about a specific project in the gitlab provider. The results include the name of the project, path, description, default branch, etc.

## Example Usage

```hcl
data "gitlab_projects" "example" {
}
```

## Argument Reference

The following arguments are supported:

* `order_by` - (Optional) Return projects ordered by id, name, path, created_at, updated_at, or last_activity_at fields. Default is created_at.

* `search` - (Optional) Return list of projects matching the search criteria

## Attributes Reference

The following attributes are exported:

* `projects` - A list containing the projects matching the supplied arguments

Projects have the following fields:

* `id` - The ID of the project

* `name` - The name of the project
