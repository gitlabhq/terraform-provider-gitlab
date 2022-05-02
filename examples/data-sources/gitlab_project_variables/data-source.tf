data "gitlab_project_variables" "vars" {
  project = "my/example/project"
}

# Using an environment scope
data "gitlab_project_variables" "staging_vars" {
  project           = "my/example/project"
  environment_scope = "staging/*"
}
