# By project ID
data "gitlab_project_tag" "foo" {
  name    = "example"
  project = "12345"
}

# By project full path
data "gitlab_project_tag" "foo" {
  name    = "example"
  project = "foo/bar"
}
