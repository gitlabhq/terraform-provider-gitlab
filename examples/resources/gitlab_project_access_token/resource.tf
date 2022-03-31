resource "gitlab_project_access_token" "example" {
  project      = "25"
  name         = "Example project access token"
  expires_at   = "2020-03-14"
  access_level = "reporter"

  scopes = ["api"]
}

resource "gitlab_project_variable" "example" {
  project = gitlab_project.example.id
  key     = "pat"
  value   = gitlab_project_access_token.example.token
}
