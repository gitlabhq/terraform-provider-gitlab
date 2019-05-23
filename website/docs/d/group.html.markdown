---
layout: "gitlab"
page_title: "GitLab: gitlab_group"
sidebar_current: "docs-gitlab-data-source-group"
description: |-
  Looks up a gitlab group
---

# gitlab\_group

Provides details about a specific group in the gitlab provider.

## Example Usage

**By group's ID**

```hcl
data "gitlab_group" "foo" {
  group_id = 123
}
```

**By group's full path**

```hcl
data "gitlab_group" "foo" {
  full_path = "foo/bar"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The ID of the group.

* `full_path` - (Optional) The full path of the group.

**Note**: exactly one of group_id or full_path must be provided.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique ID assigned to the group.

* `name` - The name of this group.

* `path` - The path of the group.

* `description` - The description of the group.

* `lfs_enabled` - Boolean, is LFS enabled for projects in this group.

* `request_access_enabled` - Boolean, is request for access enabled to the group.

* `visibility_level` - Visibility level of the group. Possible values are `private`, `internal`, `public`.

* `parent_id` - Integer, ID of the parent group.
  
* `full_path` - The full path of the group.

* `full_name` - The full name of the group.

* `web_url` - Web URL of the group.

[doc]: https://docs.gitlab.com/ee/api/groups.html#details-of-a-group
