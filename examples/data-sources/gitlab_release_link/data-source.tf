# By project ID
data "gitlab_release_link" "example" {
  project  = "12345"
  tag_name = "v1.0.1"
  link_id  = "11"
}

# By project full path
data "gitlab_release_link" "example" {
  project  = "foo/bar"
  tag_name = "v1.0.1"
  link_id  = "11"
}
