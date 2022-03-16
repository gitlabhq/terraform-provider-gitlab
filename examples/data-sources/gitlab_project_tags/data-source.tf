# By project ID
data "gitlab_project_tags" "example" {
  project = "12345"
}

# By project full path
data "gitlab_project_tags" "example" {
  project = "foo/bar"
}
