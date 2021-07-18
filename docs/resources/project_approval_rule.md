# gitlab\_project\_approval\_rule

This resource allows you to create and manage multiple approval rules for your GitLab
projects. For further information on approval rules, consult the [gitlab
documentation](https://docs.gitlab.com/ee/api/merge_request_approvals.html#project-level-mr-approvals).

-> This feature requires GitLab Premium.

## Example Usage

```hcl
resource "gitlab_project_approval_rule" "example-one" {
  project            = 5
  name               = "Example Rule"
  approvals_required = 3
  user_ids           = [50, 500]
  group_ids          = [51]
}
```

### With Protected Branch IDs

```hcl
resource "gitlab_branch_protection" "example" {
  project            = 5
  branch             = "release/*"
  push_access_level  = "maintainer"
  merge_access_level = "developer"
}

resource "gitlab_project_approval_rule" "example" {
  project              = 5
  name                 = "Example Rule"
  approvals_required   = 3
  user_ids             = [50, 500]
  group_ids            = [51]
  protected_branch_ids = [gitlab_branch_protection.example.branch_protection_id]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the approval rules.

* `name` - (Required) The name of the approval rule.

* `approvals_required` - (Required) The number of approvals required for this rule.

* `user_ids` - (Optional)  A list of specific User IDs to add to the list of approvers.

* `group_ids` - (Optional) A list of group IDs whose members can approve of the merge request.

* `protected_branch_ids` - (Optional) A list of protected branch IDs (not branch names) for which the rule applies.

## Import

GitLab project approval rules can be imported using an id consisting of `project-id:rule-id`, e.g.

```
$ terraform import gitlab_project_approval_rule.example "12345:6"
```
