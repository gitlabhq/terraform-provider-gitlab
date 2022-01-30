# A GitLab Group Access Token can be imported using a key composed of `<group-id>:<token-id>`, e.g.
terraform import gitlab_group_access_token.example "12345:1"

# ATTENTION: the `token` resource attribute is not available for imported resources as this information cannot be read from the GitLab API.
