---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_group_membership Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_group_membership resource allows to manage the lifecycle of a users group membership.
  -> If a group should grant membership to another group use the gitlab_group_share_group resource instead.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/members/
---

# gitlab_group_membership (Resource)

The `gitlab_group_membership` resource allows to manage the lifecycle of a users group membership.

-> If a group should grant membership to another group use the `gitlab_group_share_group` resource instead.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/members/)

## Example Usage

```terraform
resource "gitlab_group_membership" "test" {
  group_id     = 12345
  user_id      = 1337
  access_level = "guest"
  expires_at   = "2020-12-31"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `access_level` (String) Access level for the member. Valid values are: `no one`, `minimal`, `guest`, `planner`, `reporter`, `developer`, `maintainer`, `owner`.
- `group_id` (Number) The ID of the group.
- `user_id` (Number) The ID of the user.

### Optional

- `expires_at` (String) Expiration date for the group membership. Format: `YYYY-MM-DD`
- `member_role_id` (Number) The ID of a custom member role. Only available for Ultimate instances.
- `skip_subresources_on_destroy` (Boolean) Whether the deletion of direct memberships of the removed member in subgroups and projects should be skipped. Only used during a destroy.
- `unassign_issuables_on_destroy` (Boolean) Whether the removed member should be unassigned from any issues or merge requests inside a given group or project. Only used during a destroy.

### Read-Only

- `id` (String) The ID of the group membership. In the format of `<group-id:user-id>`.

## Import

Starting in Terraform v1.5.0, you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_group_membership`. For example:

```terraform
import {
  to = gitlab_group_membership.example
  id = "see CLI command below for ID"
}
```

Importing using the CLI is supported with the following syntax:

```shell
# GitLab group membership can be imported using an id made up of `group_id:user_id`, e.g.
terraform import gitlab_group_membership.test "12345:1337"
```
