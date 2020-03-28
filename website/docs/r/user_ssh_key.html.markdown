---
layout: "gitlab"
page_title: "GitLab: gitlab_user_ssh_key"
sidebar_current: "docs-gitlab-resource-user-ssh-key"
description: |-
  Provides a Gitlab user's SSH key resource.
---

# gitlab_user_ssh_key

Provides a Gitlab user's SSH key resource.

This resource allows you to add/remove SSH keys from your user account.

## Example Usage

```hcl
resource "gitlab_user_ssh_key" "example" {
  title = "example title"
  key   = "${file("~/.ssh/id_rsa.pub")}"
}
```

## Argument Reference

The following arguments are supported:

* `title` - (Required) A descriptive name for the new key. e.g. `Personal Linux Ubuntu`
* `key` - (Required) The public SSH key to add to your Gitlab account.

## Attributes Reference

The following attributes are exported:

* `key_id` - The ID of the SSH key
