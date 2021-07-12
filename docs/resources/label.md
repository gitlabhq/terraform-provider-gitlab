# gitlab\_label

This resource allows you to create and manage labels for your GitLab projects.
For further information on labels, consult the [gitlab
documentation](https://docs.gitlab.com/ee/user/project/labels.html#project-labels).

## Example Usage

```hcl
resource "gitlab_label" "fixme" {
  project     = "example"
  name        = "fixme"
  description = "issue with failing tests"
  color       = "#ffcc00"
}

# Scoped label
resource "gitlab_label" "devops_create" {
  project     = gitlab_project.example.id
  name        = "devops::create"
  description = "issue for creating infrastructure resources"
  color       = "#ffa500"
}

```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name or id of the project to add the label to.

* `name` - (Required) The name of the label.

* `color` - (Required) The color of the label given in 6-digit hex notation with leading '#' sign (e.g. #FFAABB) or one of the [CSS color names](https://developer.mozilla.org/en-US/docs/Web/CSS/color_value#Color_keywords).

* `description` - (Optional) The description of the label.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the label by the GitLab server (the name of the label).
