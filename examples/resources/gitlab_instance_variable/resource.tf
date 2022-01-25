resource "gitlab_instance_variable" "example" {
   key       = "instance_variable_key"
   value     = "instance_variable_value"
   protected = false
   masked    = false
}
