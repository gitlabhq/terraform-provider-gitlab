data "gitlab_project" "example" {
  id = "foo/bar/baz"
}

data "gitlab_project_hook" "example" {
  project = data.gitlab_project.example.id
  hook_id = 1
}
