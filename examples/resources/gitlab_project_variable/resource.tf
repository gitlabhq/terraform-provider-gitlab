resource "gitlab_project_variable" "example" {
   project   = "12345"
   key       = "project_variable_key"
   value     = "project_variable_value"
   protected = false
}
