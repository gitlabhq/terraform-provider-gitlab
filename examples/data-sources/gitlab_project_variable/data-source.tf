data "gitlab_project_variable" "foo" {
  project = "my/example/project"
  key     = "foo"
}

# Using an environment scope
data "gitlab_project_variable" "bar" {
  project           = "my/example/project"
  key               = "bar"
  environment_scope = "staging/*"
}
