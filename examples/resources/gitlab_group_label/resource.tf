resource "gitlab_group_label" "fixme" {
  group       = "example"
  name        = "fixme"
  description = "issue with failing tests"
  color       = "#ffcc00"
}
