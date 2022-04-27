# A GitLab Personal Access Token can be imported using a key composed of `<user-id>:<token-id>`, e.g.
terraform import gitlab_personal_access_token.example "12345:1"

# NOTE: the `token` resource attribute is not available for imported resources as this information cannot be read from the GitLab API.
