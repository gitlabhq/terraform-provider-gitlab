---
layout: "gitlab"
page_title: "GitLab: gitlab_group"
sidebar_current: "docs-gitlab-resource-group"
description: |-
  Creates and manages GitLab groups
---

# gitlab\_group

This resource allows you to create and manage GitLab groups.
Note your provider will need to be configured with admin-level access for this resource to work.

## Example Usage

```hcl
resource "gitlab_group" "example" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

// Create a project in the example group
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = "${gitlab_group.example.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of this group.

* `path` - (Required) The path of the group.

* `description` - (Optional) The description of the group.

* `lfs_enabled` - (Optional) Boolean, defaults to true.  Whether to enable LFS
support for projects in this group.

* `request_access_enabled` - (Optional) Boolean, defaults to false.  Whether to
enable users to request access to the group.

* `visibility_level` - (Optional) Set to `public` to create a public group.
  Valid values are `private`, `internal`, `public`.
  Groups are created as private by default.

* `parent_id` - (Optional) Integer, id of the parent group (creates a nested group).

* `subgroup_creation_level` - (Optional) Define who can create sub groups.
  Valid values are `owner` or `maintainer`.

* `project_creation_level` - (Optional) Determine if developers can create projects in the group.
  Valid values are `noone`, `developer` and `maintainer`.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the group by the GitLab server.  Serves as a
  namespace id where one is needed.

* `full_path` - The full path of the group.

* `full_name` - The full name of the group.

* `web_url` - Web URL of the group.

* `runners_token` - The group level registration token to use during runner setup.

## Importing groups

You can import a group state using `terraform import <resource> <id>`.  The
`id` can be whatever the [details of a group][details_of_a_group] api takes for
its `:id` value, so for example:

    terraform import gitlab_group.example example

[details_of_a_group]: https://docs.gitlab.com/ee/api/groups.html#details-of-a-group
