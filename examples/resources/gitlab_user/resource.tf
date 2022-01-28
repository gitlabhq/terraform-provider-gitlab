resource "gitlab_user" "example" {
  name             = "Example Foo"
  username         = "example"
  password         = "superPassword"
  email            = "gitlab@user.create"
  is_admin         = true
  projects_limit   = 4
  can_create_group = false
  is_external      = true
  reset_password   = false
}
