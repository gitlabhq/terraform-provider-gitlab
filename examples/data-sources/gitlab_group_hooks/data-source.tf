data "gitlab_group" "example" {
  id = "foo/bar/baz"
}

data "gitlab_group_hooks" "examples" {
  group = data.gitlab_group.example.id
}
