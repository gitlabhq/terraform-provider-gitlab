---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_group Data Source - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_group data source allows details of a group to be retrieved by its id or full path.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/groups/#get-a-single-group
---

# gitlab_group (Data Source)

The `gitlab_group` data source allows details of a group to be retrieved by its id or full path.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/groups/#get-a-single-group)

## Example Usage

```terraform
# By group's ID
data "gitlab_group" "foo" {
  group_id = 123
}

# By group's full path
data "gitlab_group" "foo" {
  full_path = "foo/bar"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- `full_path` (String) The full path of the group.
- `group_id` (Number) The ID of the group.

### Read-Only

- `default_branch` (String) The default branch of the group.
- `default_branch_protection` (Number) Whether developers and maintainers can push to the applicable default branch.
- `description` (String) The description of the group.
- `extra_shared_runners_minutes_limit` (Number) Can be set by administrators only. Additional CI/CD minutes for this group.
- `full_name` (String) The full name of the group.
- `id` (String) The ID of this resource.
- `lfs_enabled` (Boolean) Boolean, is LFS enabled for projects in this group.
- `membership_lock` (Boolean) Users cannot be added to projects in this group.
- `name` (String) The name of this group.
- `parent_id` (Number) Integer, ID of the parent group.
- `path` (String) The path of the group.
- `prevent_forking_outside_group` (Boolean) When enabled, users can not fork projects from this group to external namespaces.
- `request_access_enabled` (Boolean) Boolean, is request for access enabled to the group.
- `runners_token` (String, Sensitive) The group level registration token to use during runner setup.
- `shared_runners_minutes_limit` (Number) Can be set by administrators only. Maximum number of monthly CI/CD minutes for this group. Can be nil (default; inherit system default), 0 (unlimited), or > 0.
- `shared_runners_setting` (String) Enable or disable shared runners for a group’s subgroups and projects. Valid values are: `enabled`, `disabled_and_overridable`, `disabled_and_unoverridable`, `disabled_with_override`.
- `shared_with_groups` (List of Object) Describes groups which have access shared to this group. (see [below for nested schema](#nestedatt--shared_with_groups))
- `visibility_level` (String) Visibility level of the group. Possible values are `private`, `internal`, `public`.
- `web_url` (String) Web URL of the group.
- `wiki_access_level` (String) The group's wiki access level. Only available on Premium and Ultimate plans. Valid values are `disabled`, `private`, `enabled`.

<a id="nestedatt--shared_with_groups"></a>
### Nested Schema for `shared_with_groups`

Read-Only:

- `expires_at` (String)
- `group_access_level` (Number)
- `group_full_path` (String)
- `group_id` (Number)
- `group_name` (String)
