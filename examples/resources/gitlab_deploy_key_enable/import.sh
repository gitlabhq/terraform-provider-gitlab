# GitLab enabled deploy keys can be imported using an id made up of `{project_id}:{deploy_key_id}`, e.g.
# `project_id` can be whatever the [get single project api][get_single_project] takes for
# its `:id` value, so for example:
terraform import gitlab_deploy_key_enable.example 12345:67890
terraform import gitlab_deploy_key_enable.example richardc/example:67890
