# gitlab\_project_membership

This resource allows you to add a current user to an existing project with a set access level.

## Example Usage

```hcl
resource "gitlab_project_membership" "test" {
  project_id   = "12345"
  user_id      = 1337
  access_level = "guest"
}

resource "gitlab_project_membership" "example" {
  project_id   = "67890"
  user_id      = 1234
  access_level = "guest"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The id of the project.

* `user_id` - (Required) The id of the user.

* `access_level` - (Required) One of five levels of access to the project.

## Import

GitLab group membership can be imported using an id made up of `group_id:user_id`, e.g.

```
$ terraform import gitlab_project_membership.test "12345:1337"
```
