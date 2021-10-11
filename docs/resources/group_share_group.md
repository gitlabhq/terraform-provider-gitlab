# gitlab\_group\_share\_group

This resource allows you to share a group with another group

## Example Usage

```hcl
resource "gitlab_group_share_group" "test" {
  group_id       = gitlab_group.foo.id
  share_group_id = gitlab_group.bar.id
  group_access   = "guest"
  expires_at     = "2099-01-01"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of the main group.

* `share_group_id` - (Required) The id of an additional group which will be shared with the main group.

* `group_access` - (Required) One of five levels of access to the group.

* `expires_at` - (Optional) Share expiration date. Format: `YYYY-MM-DD`

## Import

GitLab group shares can be imported using an id made up of `mainGroupId:shareGroupId`, e.g.

```
$ terraform import gitlab_group_share_group.test 12345:1337
```
