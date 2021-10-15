# gitlab\_group\_custom\_attribute

This resource allows you to set custom attributes for a group.

## Example Usage

```hcl
resource "gitlab_group_custom_attribute" "attr" {
  group = "42"
  key   = "location"
  value = "Greenland"
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The id of the group.

* `key` - (Required) Key for the Custom Attribute.

* `value` - (Required) Value for the Custom Attribute.

## Import

You can import a group custom attribute using the following id pattern:

```shell
$ terraform import gitlab_group_custom_attribute.attr <group-id>:<key>
```

For the example above this would be:

```shell
$ terraform import gitlab_group_custom_attribute.attr 42:location
```
