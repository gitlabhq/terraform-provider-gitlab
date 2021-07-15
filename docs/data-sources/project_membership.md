# gitlab\_project\_membership

Provides details about a list of project members in the gitlab provider. The results include id, username, name and more about the requested members.

## Example Usage

**By project's ID**

```hcl
data "gitlab_project_membership" "example" {
  id = 123
}
```

**By project's full path**

```hcl
data "gitlab_group_membership" "example" {
  id = "foo/bar/baz"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The integer or path with namespace that uniquely identifies the project within the gitlab install.

* `access_level` - (Optional) Only return members with the desired access level. Acceptable values are: `guest`, `reporter`, `developer`, `maintainer`, `owner`.

* `inherited` - (Optional) Return all members of a project, even those that are members through ancestor groups. Defaults to `true`.

## Attributes Reference

The following attributes are exported:

* `members` - The list of group members.
  * `id` - The unique id assigned to the user by the gitlab server.
  * `username` - The username of the user.
  * `name` - The name of the user.
  * `state` - Whether the user is active or blocked.
  * `avatar_url` - The avatar URL of the user.
  * `web_url` - User's website URL.
  * `access_level` - One of five levels of access to the project.
  * `expires_at` - Expiration date for the project membership.
