# By project ID
data "gitlab_branch" "foo" {
  name    = "example"
  project = "12345"
}

# By project full path
data "gitlab_branch" "foo" {
  name    = "example"
  project = "foo/bar"
}