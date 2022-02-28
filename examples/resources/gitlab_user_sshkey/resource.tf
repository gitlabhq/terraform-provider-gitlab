data "gitlab_user" "example" {
  username = "example-user"
}

resource "gitlab_user_sshkey" "example" {
  user_id    = data.gitlab_user.id
  title      = "example-key"
  key        = "ssh-rsa AAAA..."
  expires_at = "2016-01-21T00:00:00.000Z"
}
