# gitlab\_tag\_protection

This resource allows you to protect a specific tag or wildcard by an access level so that the user with less access level cannot Create the tags.

## Example Usage

```hcl
resource "gitlab_tag_protection" "TagProtect" {
  project = "12345"
  tag = "TagProtected"
  create_access_level = "developer"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project.

* `tag` - (Required) Name of the tag or wildcard.

* `create_access_level` - (Required) One of five levels of access to the project.

## Importing tag protections

You can import a tag protection state using `terraform import <resource> <project_id>:<tag_name>`.

    terraform import gitlab_tag_protection.example 123456789:v1.0.0
