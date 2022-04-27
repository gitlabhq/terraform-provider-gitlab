resource "gitlab_personal_access_token" "example" {
  user_id    = "25"
  name       = "Example personal access token"
  expires_at = "2020-03-14"

  scopes = ["api"]
}

resource "gitlab_project_variable" "example" {
  project = gitlab_project.example.id
  key     = "pat"
  value   = gitlab_personal_access_token.example.token
}
