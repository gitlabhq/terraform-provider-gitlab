resource "gitlab_group_hook" "example" {
  group                 = "example/hooked"
  url                   = "https://example.com/hook/example"
  merge_requests_events = true
}

# Setting all attributes
resource "gitlab_group_hook" "all_attributes" {
  group                      = 1
  url                        = "http://example.com"
  token                      = "supersecret"
  enable_ssl_verification    = false
  push_events                = true
  push_events_branch_filter  = "devel"
  issues_events              = false
  confidential_issues_events = false
  merge_requests_events      = true
  tag_push_events            = true
  note_events                = true
  confidential_note_events   = true
  job_events                 = true
  pipeline_events            = true
  wiki_page_events           = true
  deployment_events          = true
  releases_events            = true
  subgroup_events            = true
}
