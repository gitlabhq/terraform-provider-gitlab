data "gitlab_project_protected_branches" "example" {
  project_id = 30
}

data "gitlab_project_protected_branches" "example" {
  project_id = "foo/bar/baz"
}
