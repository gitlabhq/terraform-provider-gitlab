data "gitlab_project" "example" {
  id = "foo/bar/baz"
}

data "gitlab_project_hooks" "examples" {
  project = data.gitlab_project.example.id
}
