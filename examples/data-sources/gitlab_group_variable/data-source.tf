data "gitlab_group_variable" "foo" {
  group = "my/example/group"
  key   = "foo"
}

# Using an environment scope
data "gitlab_group_variable" "bar" {
  group             = "my/example/group"
  key               = "bar"
  environment_scope = "staging/*"
}
