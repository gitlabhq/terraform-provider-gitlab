---
layout: "gitlab"
page_title: "GitLab: gitlab_branch_protection"
sidebar_current: "docs-gitlab-resource-branch_protection"
description: |-
  Protects a branch by assigning access levels to it
---

# gitlab\_branch_protection

This resource allows you to protect a specific branch by an access level so that the user with less access level cannot Merge/Push to the branch. GitLab EE features to protect by group or user are also supported.

## Example Usage

```hcl
resource "gitlab_branch_protection" "BranchProtect" {
  project = "12345"
  branch = "BranchProtected"
  push_access_level = "developer"
  merge_access_level = "developer"
}
```

```hcl
resource "gitlab_branch_protection" "BranchProtect" {
  project = "12345"
  branch = "BranchProtected"
  allowed_to_push {
    user_id = [123, 124]
    group_id = [1, 2]
    access_level = ["developer"]
  }
  allowed_to_merge {
    access_level = ["maintainer"]
  }
  code_owner_approval_required = true
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project.

* `branch` - (Required) Name of the branch.

* `push_access_level` - (Required if not setting allowed_to_push) One of five levels of access to the project.

* `merge_access_level` - (Required if not setting allowed_to_merge) One of five levels of access to the project.

* `unprotect_access_level` - (Optional; conflicts with allowed_to_unprotect) One of five levels of access to the project.  Defaults to Maintainer.

* `allowed_to_push` - (Required if not setting push_access_level) GitLab EE Only - Map of user ids, group ids, and access levels to grant access.  At least one of user_id, group_id, or access_level must be defined.

* `allowed_to_merge` - (Required if not setting merge_access_level) GitLab EE Only - Map of user ids, group ids, and access levels to grant access.  At least one of user_id, group_id, or access_level must be defined.

* `allowed_to_unprotect` - (Optional; conflicts with unprotect_access_level) GitLab EE Only - Map of user ids, group ids, and access levels to grant access.  At least one of user_id, group_id, or access_level must be defined.

* `code_owner_approval_required` - (Optional) GitLab EE Only - True or false.  Defaults to false.