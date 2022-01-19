# gitlab\_topic

This resource allows you to create and manage topics that are then assignable to projects. Topics are the successors for project tags. Aside from avoiding terminology collisions with Git tags, they are more descriptive and better searchable.

For assigning topics, use the [project](./project.md) resource.

## Example Usage

```hcl
resource "gitlab_topic" "functional-programming" {
  name             = "Functional Programming"
  description      = "In computer science, functional programming is a programming paradigm where programs are constructed by applying and composing functions."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the approval rule.

* `description` - (Optional) A text describing the topic.

## Import

GitLab topics can also be imported, e.g.

```
$ terraform import gitlab_topic.functional-programming
```
