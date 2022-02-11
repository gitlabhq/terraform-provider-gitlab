resource "gitlab_group_membership" "test" {
  group_id     = "12345"
  user_id      = 1337
  access_level = "guest"
  expires_at   = "2020-12-31"
}
