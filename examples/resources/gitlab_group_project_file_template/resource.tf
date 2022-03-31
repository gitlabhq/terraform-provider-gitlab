resource "gitlab_group" "foo" {
  name        = "group"
  path        = "group"
  description = "An example group"
}

resource "gitlab_project" "bar" {
  name             = "template project"
  description      = "contains file templates"
  visibility_level = "public"

  namespace_id = gitlab_group.foo.id
}

resource "gitlab_group_project_file_template" "template_link" {
  group_id = gitlab_group.foo.id
  project  = gitlab_project.bar.id
}
