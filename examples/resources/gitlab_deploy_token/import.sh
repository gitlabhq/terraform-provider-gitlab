# GitLab deploy tokens can be imported using an id made up of `{type}:{type_id}:{deploy_token_id}`, where type is one of: project, group.
terraform import gitlab_deploy_token.group_token group:1:3
terraform import gitlab_deploy_token.project_token project:1:4

# Note: the `token` resource attribute is not available for imported resources as this information cannot be read from the GitLab API.
