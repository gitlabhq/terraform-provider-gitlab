# GitLab deploy keys can be imported using an id made up of `{project_id}:{deploy_key_id}`, e.g.
# `project_id` can be whatever the [get single project api][get_single_project] takes for
# its `:id` value, so for example:
terraform import gitlab_deploy_key.test 1:3
terraform import gitlab_deploy_key.test richardc/example:3
