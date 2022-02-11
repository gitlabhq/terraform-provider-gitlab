resource "gitlab_project_hook" "example" {
  project               = "example/hooked"
  url                   = "https://example.com/hook/example"
  merge_requests_events = true
}
