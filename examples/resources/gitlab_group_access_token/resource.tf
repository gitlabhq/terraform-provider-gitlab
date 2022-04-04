resource "gitlab_group_access_token" "example" {
  group        = "25"
  name         = "Example project access token"
  expires_at   = "2020-03-14"
  access_level = "developer"

  scopes = ["api"]
}

resource "gitlab_group_variable" "example" {
  group = "25"
  key   = "gat"
  value = gitlab_group_access_token.example.token
}
