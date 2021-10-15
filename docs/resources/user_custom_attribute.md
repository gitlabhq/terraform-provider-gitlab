# gitlab\_user\_custom\_attribute

This resource allows you to set custom attributes for a user.

## Example Usage

```hcl
resource "gitlab_user_custom_attribute" "attr" {
  user  = "42"
  key   = "location"
  value = "Greenland"
}
```

## Argument Reference

The following arguments are supported:

* `user` - (Required) The id of the user.

* `key` - (Required) Key for the Custom Attribute.

* `value` - (Required) Value for the Custom Attribute.

## Import

You can import a user custom attribute using the following id pattern:

```shell
$ terraform import gitlab_user_custom_attribute.attr <user-id>:<key>
```

For the example above this would be:

```shell
$ terraform import gitlab_user_custom_attribute.attr 42:location
```
