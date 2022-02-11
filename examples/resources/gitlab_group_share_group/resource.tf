resource "gitlab_group_share_group" "test" {
  group_id       = gitlab_group.foo.id
  share_group_id = gitlab_group.bar.id
  group_access   = "guest"
  expires_at     = "2099-01-01"
}
