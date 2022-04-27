# By project ID
data "gitlab_project_milestones" "example" {
  project_id = "12345"
}

# By project full path
data "gitlab_project_milestones" "example" {
  project_id = "foo/bar"
}
