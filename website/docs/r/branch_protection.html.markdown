---
layout: "gitlab"
page_title: "GitLab: gitlab_branch_protection"
sidebar_current: "docs-gitlab-resource-branch_protection"
description: |-
  Protects a branch by assigning access levels to it
---

# gitlab\_branch_protection

This resource allows you to protect a specific branch by an access level so that the user with less access level cannot Merge/Push to the branch. GitLab EE features to protect by group or user are not supported.

## Example Usage

```hcl
resource "gitlab_branch_protection" "BranchProtect" {
  project = "12345"
  branch = "BranchProtected"
  push_access_level = "developer"
  merge_access_level = "developer"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project.

* `branch` - (Required) Name of the branch.

* `push_access_level` - (Required) One of five levels of access to the project.

* `merge_access_level` - (Required) One of five levels of access to the project.