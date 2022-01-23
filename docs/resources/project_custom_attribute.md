# gitlab\_project\_custom\_attribute

This resource allows you to set custom attributes for a project.

## Example Usage

```hcl
resource "gitlab_project_custom_attribute" "attr" {
  project = "42"
  key     = "location"
  value   = "Greenland"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project.

* `key` - (Required) Key for the Custom Attribute.

* `value` - (Required) Value for the Custom Attribute.

## Import

You can import a project custom attribute using the following id pattern:

```shell
$ terraform import gitlab_project_custom_attribute.attr <project-id>:<key>
```

For the example above this would be:

```shell
$ terraform import gitlab_project_custom_attribute.attr 42:location
```
