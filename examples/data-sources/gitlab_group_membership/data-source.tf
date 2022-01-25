# By group's ID
data "gitlab_group_membership" "example" {
  group_id = 123
}

# By group's full path
data "gitlab_group_membership" "example" {
  full_path = "foo/bar"
}
