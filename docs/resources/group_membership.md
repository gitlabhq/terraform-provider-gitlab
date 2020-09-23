# gitlab\_group_membership

This resource allows you to add a user to an existing group.

## Example Usage

```hcl
resource "gitlab_group_membership" "test" {
  group_id     = "12345"
  user_id      = 1337
  access_level = "guest"
  expires_at   = "2020-12-31"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of the group.

* `user_id` - (Required) The id of the user.

* `access_level` - (Required)  Acceptable values are: guest, reporter, developer, maintainer, owner.

* `expires_at` - (Optional) Expiration date for the group membership. Format: `YYYY-MM-DD`

## Import

GitLab group membership can be imported using an id made up of `group_id:user_id`, e.g.

```
$ terraform import gitlab_group_membership.test "12345:1337"
```
