# A token for a GitLab Agent for Kubernetes can be imported with the following command and the id pattern `<project>:<agent-id>:<token-id>`:
terraform import gitlab_cluster_agent_token.example '12345:42:1'

# ATTENTION: the `token` resource attribute is not available for imported resources as this information cannot be read from the GitLab API.
