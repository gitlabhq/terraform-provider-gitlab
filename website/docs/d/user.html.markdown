---
layout: "gitlab"
page_title: "GitLab: gitlab_user"
sidebar_current: "docs-gitlab-data-source-user"
description: |-
  Looks up a gitlab user
---

# gitlab\_user

Provides details about a specific user in the gitlab provider. The results include username, id, name, etc.

## Example Usage

```hcl
data "gitlab_user" "example" {
	email = "test@aaa.com"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The e-mail address of the user.

## Attributes Reference

The following attributes are exported:

* `name` - The name of the user.

* `username` - The username of the user.

* `email` - The e-mail address of the user.

* `id` - The unique id assigned to the user by the gitlab server.


