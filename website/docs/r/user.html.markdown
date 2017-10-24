---
layout: "gitlab"
page_title: "GitLab: gitlab_user"
sidebar_current: "docs-gitlab-resource-user"
description: |-
  Creates and manages GitLab users
---

# gitlab\_user

This resource allows you to create and manage GitLab users.
Note your provider will need to be configured with admin-level access for this resource to work.

## Example Usage

```hcl
resource "gitlab_user" "example" {
  name             = "Example Foo"
  username         = "example"
  password         = "superPassword"
  email            = "gitlab@user.create"
  is_admin         = true
  projects_limit   = 4
  can_create_group = false
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the user.

* `username` - (Required) The username of the user.

* `password` - (Required) The password of the user.

* `email` - (Required) The e-mail address of the user.

* `is_admin` - (Optional) Boolean, defaults to false.  Whether to enable administrative priviledges
for the user.

* `projects_limit` - (Optional) Integer, defaults to 0.  Number of projects user can create.

* `can_create_group` - (Optional) Boolean, defaults to false. Whether to allow the user to create groups.

* `skip_confirmation` - (Optional) Boolean, defaults to true. Whether to skip confirmation.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the user by the GitLab server.

## Importing users

You can import a user to terraform state using `terraform import <resource> <id>`.
The `id` must be an integer for the id of the user you want to import,
for example:

    terraform import gitlab_user.example 42
