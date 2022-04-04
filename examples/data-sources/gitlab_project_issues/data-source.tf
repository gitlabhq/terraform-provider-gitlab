data "gitlab_project" "foo" {
  id = "foo/bar/baz"
}

data "gitlab_project_issues" "all_with_foo" {
  project = data.gitlab_project.foo.id
  search  = "foo"
}
