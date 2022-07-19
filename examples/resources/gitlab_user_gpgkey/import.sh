# You can import a GPG key for a specific user using an id made up of `{user-id}:{key}`, e.g.
terraform import gitlab_user_gpgkey.example 42:1

# Alternatively, you can import a GPG key for the current user using an id made up of `{key}`, e.g.
terraform import gitlab_user_gpgkey.example_user 1
