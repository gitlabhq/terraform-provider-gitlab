# gitlab\_group\_membership

Provide details about a list of group members in the gitlab provider. The results include id, username, name and more about the requested members.

## Example Usage

### By group's ID

```hcl
data "gitlab_group_membership" "example" {
  group_id = 123
}
```

### By group's full path

```hcl
data "gitlab_group_membership" "example" {
  full_path = "foo/bar"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Optional) The ID of the group.

* `full_path` - (Optional) The full path of the group.

* `access_level` - (Optional) Only return members with the desired access level. Acceptable values are: `guest`, `reporter`, `developer`, `maintainer`, `owner`.

> **Note**: exactly one of group_id or full_path must be provided.

## Attributes Reference

The following attributes are exported:

* `members` - The list of group members.
  * `id` - The unique id assigned to the user by the gitlab server.
  * `username` - The username of the user.
  * `name` - The name of the user.
  * `state` - Whether the user is active or blocked.
  * `avatar_url` - The avatar URL of the user.
  * `web_url` - User's website URL.
  * `access_level` - One of five levels of access to the group.
  * `expires_at` - Expiration date for the group membership.
