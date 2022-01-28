data "gitlab_project_protected_branch" "example" {
  project_id = 30
  name       = "main"
}

data "gitlab_project_protected_branch" "example" {
  project_id = "foo/bar/baz"
  name       = "main"
}
