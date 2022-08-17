# Set the default branch
resource "gitlab_application_settings" "this" {
  default_branch_name = "main"
}

# Set the 2FA settings
resource "gitlab_application_settings" "this" {
  require_two_factor_authentication = true
  two_factor_grace_period           = 24
}
