# gitlab\_project\_share\_group

This resource allows you to share a project with a group

## Example Usage

```hcl
resource "gitlab_project_share_group" "test" {
  project_id = "12345"
  group_id = 1337
  access_level = "guest"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The id of the project.

* `group_id` - (Required) The id of the group.

* `access_level` - (Required) One of five levels of access to the project.

## Import

GitLab project group shares can be imported using an id made up of `projectid:groupid`, e.g.

```
$ terraform import gitlab_project_share_group.test 12345:1337
```
