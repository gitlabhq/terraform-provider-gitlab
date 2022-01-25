resource "gitlab_project_share_group" "test" {
  project_id = "12345"
  group_id = 1337
  access_level = "guest"
}
