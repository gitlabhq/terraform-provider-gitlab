# By project ID
data "gitlab_release_links" "example" {
  project  = "12345"
  tag_name = "v1.0.1"
}

# By project full path
data "gitlab_release_links" "example" {
  project  = "foo/bar"
  tag_name = "v1.0.1"
}
