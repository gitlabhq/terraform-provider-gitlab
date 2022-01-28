resource "gitlab_group_variable" "example" {
  group             = "12345"
  key               = "group_variable_key"
  value             = "group_variable_value"
  protected         = false
  masked            = false
  environment_scope = "*"
}
