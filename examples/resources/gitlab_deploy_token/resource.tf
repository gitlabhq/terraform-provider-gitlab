# Example Usage - Project
resource "gitlab_deploy_token" "example" {
  project    = "example/deploying"
  name       = "Example deploy token"
  username   = "example-username"
  expires_at = "2020-03-14T00:00:00.000Z"

  scopes = ["read_repository", "read_registry"]
}

resource "gitlab_deploy_token" "example-two" {
  project    = "12345678"
  name       = "Example deploy token expires in 24h"
  expires_at = timeadd(timestamp(), "24h")
}

# Example Usage - Group
resource "gitlab_deploy_token" "example" {
  group = "example/deploying"
  name  = "Example group deploy token"

  scopes = ["read_repository"]
}
