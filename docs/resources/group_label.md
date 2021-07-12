# gitlab\_group\_label

This resource allows you to create and manage labels for your GitLab groups.
For further information on labels, consult the [gitlab
documentation](https://docs.gitlab.com/ee/user/project/labels.html#group-labels).

## Example Usage

```hcl
resource "gitlab_group_label" "fixme" {
  group       = "example"
  name        = "fixme"
  description = "issue with failing tests"
  color       = "#ffcc00"
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required) The name or id of the group to add the label to.

* `name` - (Required) The name of the label.

* `color` - (Required) The color of the label given in 6-digit hex notation with leading '#' sign (e.g. #FFAABB) or one of the [CSS color names](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value#Color_keywords).

* `description` - (Optional) The description of the label.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the label by the GitLab server (the name of the label).

## Import

Gitlab group labels can be imported using an id made up of `{group_id}:{group_label_id}`, e.g.

```
$ terraform import gitlab_group_label.example 12345:fixme
```
