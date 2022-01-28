# By group's ID
data "gitlab_group" "foo" {
  group_id = 123
}

# By group's full path
data "gitlab_group" "foo" {
  full_path = "foo/bar"
}
