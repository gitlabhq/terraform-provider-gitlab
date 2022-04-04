resource "gitlab_system_hook" "example" {
  url                      = "https://example.com/hook-%d"
  token                    = "secret-token"
  push_events              = true
  tag_push_events          = true
  merge_requests_events    = true
  repository_update_events = true
  enable_ssl_verification  = true
}
