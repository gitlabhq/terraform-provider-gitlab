data "gitlab_group_variables" "vars" {
  group = "my/example/group"
}

# Using an environment scope
data "gitlab_group_variables" "staging_vars" {
  group             = "my/example/group"
  environment_scope = "staging/*"
}
