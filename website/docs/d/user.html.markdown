---
layout: "gitlab"
page_title: "GitLab: gitlab_user"
sidebar_current: "docs-gitlab-data-source-user"
description: |-
  Looks up a gitlab user
---

# gitlab\_user

Provides details about a specific user in the gitlab provider. Especially the ability to lookup the id for linking to other resources.

## Example Usage

```hcl
data "gitlab_user" "example" {
	email = "test@aaa.com"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Optional) The e-mail address of the user. (Requires administrator privileges)

* `username` - (Optional) The username of the user.

If both are given only e-mail is used.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the user.

* `username` - The username of the user.

* `email` - The e-mail address of the user.

* `id` - The unique id assigned to the user by the gitlab server.


