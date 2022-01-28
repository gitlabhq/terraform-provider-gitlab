# gitlab\_managed\_license

This resource allows you to create and manage license rules for GitLab's
[License Compliance](https://docs.gitlab.com/ee/user/compliance/license_compliance/)

-> This features requires GitLab Ultimate

## Example Usage - Project

```hcl
resource "gitlab_managed_license" "license" {
  project = "example-project"
  name = "MIT"
  approval_status = "approved"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add
the managed license to.

* `name` - (Required, string) The name of the license.

* `approval_status` - (Required, string) The approval status of the specified license.
Must be either `approved` or `blacklisted`

## Attributes Reference

There are no additional attributes exported.

## Import

Managed Licenses can be imported using an ID made up of the `project_id:license_id`
combination, e.g.
```
$ terraform import gitlab_managed_license.foo "1234:5678"
```
