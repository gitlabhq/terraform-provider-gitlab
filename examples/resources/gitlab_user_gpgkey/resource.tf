data "gitlab_user" "example" {
  username = "example-user"
}

# Manages a GPG key for the specified user. An admin token is required if `user_id` is specified.
resource "gitlab_user_gpgkey" "example" {
  user_id = data.gitlab_user.example.id
  key     = "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----"
}

# Manages a GPG key for the current user
resource "gitlab_user_gpgkey" "example_user" {
  key = "-----BEGIN PGP PUBLIC KEY BLOCK-----\n...\n-----END PGP PUBLIC KEY BLOCK-----"
}
