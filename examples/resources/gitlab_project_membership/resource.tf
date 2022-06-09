resource "gitlab_project_membership" "test" {
  project_id   = "12345"
  user_id      = 1337
  access_level = "guest"
}

resource "gitlab_project_membership" "example" {
  project_id   = "67890"
  user_id      = 1234
  access_level = "guest"
  expires_at   = "2022-12-31"
}
