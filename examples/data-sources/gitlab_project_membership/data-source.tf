# By project's ID
data "gitlab_project_membership" "example" {
  project_id = 123
}

# By project's full path
data "gitlab_project_membership" "example" {
  full_path = "foo/bar"
}

# Get members of a project including all members
# through ancestor groups
data "gitlab_project_membership" "example" {
  project_id = 123
  inherited  = true
}
