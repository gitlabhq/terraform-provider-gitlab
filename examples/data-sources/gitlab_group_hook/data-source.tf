data "gitlab_group" "example" {
  id = "foo/bar/baz"
}

data "gitlab_group_hook" "example" {
  group   = data.gitlab_group.example.id
  hook_id = 1
}
