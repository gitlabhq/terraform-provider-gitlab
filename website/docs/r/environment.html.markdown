---
layout: "gitlab"
page_title: "Gitlab: gitlab_environment_file"
sidebar_current: "docs-gitlab-environment-file"
---

# gitlab\_environment

This resource allows you to create and manage Gitlab environments.

## Example Usage

```hcl
resource "gitlab_group" "this" {
    name = "example"
    path = "example"
    description = "An example project"
}

resource "gitlab_project" "this" {
    name = "exmaple"
    namespace_id = gitlab_group.this.id
    pipelines_enabled = true
}

resource "gitlab_environment" "this" {
    project = gitlab_project.this.id
    name = "dev"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The ID of the project the environment belongs to.

* `name` - (Required) The name of the environment.

* `external_url` - (Optional) An external URL to refer to.

## Attribute Reference

The resource exports the following attributes:

* `id` - The unique ID assigned to the environment.

* `slug` - The slug of the environment name.

* `state` - The state of the environment.
