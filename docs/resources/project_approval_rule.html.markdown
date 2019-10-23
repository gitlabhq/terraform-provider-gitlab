---
layout: "gitlab"
page_title: "GitLab: gitlab_project_approval_rule"
sidebar_current: "docs-gitlab-resource-project-approval-rule"
description: |-
  Creates and manages project-level approval rules for GitLab projects
---

# gitlab\_project\_approval\_rule

This resource allows you to create and manage approval rules for your GitLab 
projects. For further information on approval rules, consult the [gitlab
documentation](https://docs.gitlab.com/ee/api/merge_request_approvals.html#project-level-mr-approvals).

*Note*: This feature requires a GitLab Starter account or above.

## Example Usage

```hcl
resource "gitlab_project_approval_rule" "example" {
  project            = 5
  name               = "Example"
  approvals_required = 3
  user_ids           = [50, 500]
  group_ids          = [51]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the approval rules.

* `name` - (Required) The name of the approval rule.

* `approvals_required` - (Required) The number of approvals required for this rule.

* `user_ids` - (Optional)  A list of specific User IDs to add to the list of approvers.

* `group_ids` - (Optional) A list of group IDs who's members can approve of the merge request

## Attributes Reference

The resource exports the following additional attributes:

* `eligible_approvers` - A list of approvers from the `user_ids` parameter, groups, and the project that are included in the approval rule.

## Import

GitLab project approval rules can be imported using an id consisting of `project-id:rule-id`, e.g.

```
$ terraform import gitlab_project_approval_rule.example "12345:6"
```
