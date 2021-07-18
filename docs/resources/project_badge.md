# gitlab\_project\_badge

This resource allows you to create and manage badges for your GitLab projects.
For further information on hooks, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/project/badges.html).

## Example Usage

```hcl
resource "gitlab_project" "foo" {
  name = "foo-project"
}

resource "gitlab_project_badge" "example" {
  project   = gitlab_project.foo.id
  link_url  = "https://example.com/badge-123"
  image_url = "https://example.com/badge-123.svg"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project to add the badge to.

* `link_url` - (Required) The url linked with the badge.

* `image_url` - (Required) The image url which will be presented on project overview.

## Attributes Reference

The resource exports the following attributes:

* `rendered_link_url` - The link_url argument rendered (in case of use of placeholders).

* `rendered_image_url` - The image_url argument rendered (in case of use of placeholders).

## Import

GitLab project badges can be imported using an id made up of `{project_id}:{badge_id}`,
 e.g.

```bash
terraform import gitlab_project_badge.foo 1:3
```
