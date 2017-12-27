---
layout: "gitlab"
page_title: "GitLab: gitlab_datasource_user"
sidebar_current: "docs-gitlab-data-source-user"
description: |-
  Looks up a gitlab user
---

# gitlab\_datasource_user

datasource_user provides details about a specific user in the gitlab provier. The results include username, id, name, etc.

## Example Usage

```hcl
data "gitlab_user" "test" {
	email = "test@aaa.com"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) The e-mail address of the user.

## Attributes Reference

The following attributes are exported:

* `name` - (Required) The name of the user.

* `username` - (Required) The username of the user.

* `email` - (Required) The e-mail address of the user.

* `id` - The unique id assigned to the user by the GitLab server.


