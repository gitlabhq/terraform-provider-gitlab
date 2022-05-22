# By project ID
data "gitlab_project_milestone" "example" {
  project      = "12345"
  milestone_id = 10
}

# By project full path
data "gitlab_project_milestone" "example" {
  project      = "foo/bar"
  milestone_id = 10
}
