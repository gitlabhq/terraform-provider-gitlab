---
page_title: "Terraform GitLab Provider Version 19.0 Upgrade Guide"
subcategory: "Upgrade Guides"
---

# Upgrade to Terraform GitLab Provider Version 19.0

This is a draft and is subject to change.

The GitLab 19.0 major milestone introduced some breaking changes that this release addresses.
The provider also has some breaking changes which may require actions on the users side.
These are described below:

## Resource renames

If you are using any of the following resources, you will need to complete these actions:

- Rename the resource reference in your terraform code. For example:

  ```hcl
  resource "gitlab_integration_jira" "jira_setup" {
    ...
  }
  ```

  Becomes

  ```hcl
  resource "gitlab_project_integration_jira" "jira_setup" {
    ...
  }
  ```

- Perform a state move using one of the terraform options:
  - Command line with [state mv](https://developer.hashicorp.com/terraform/cli/commands/state/mv). For example:

  ```bash
  terraform state mv 'gitlab_integration_jira.jira_setup' 'gitlab_project_integration_jira.jira_setup'
  ```

  - State [moved block](https://developer.hashicorp.com/terraform/language/block/moved) and running an apply. For example:

  ```hcl
  moved {
    from = gitlab_integration_jira.jira_setup
    to   = gitlab_project_integration_jira.jira_setup
  }
  ```

### Renamed resources

- `gitlab_deploy_token` renamed to `gitlab_project_deploy_token`
- `gitlab_integration_custom_issue_tracker` renamed to `gitlab_project_integration_custom_issue_tracker`
- `gitlab_integration_emails_on_push` renamed to `gitlab_project_integration_emails_on_push`
- `gitlab_integration_external_wiki_resource` renamed to `gitlab_project_integration_external_wiki_resource`
- `gitlab_integration_github` renamed to `gitlab_project_integration_github`
- `gitlab_integration_harbor` renamed to `gitlab_project_integration_harbor`
- `gitlab_integration_jenkins` renamed to `gitlab_project_integration_jenkins`
- `gitlab_integration_jira` renamed to `gitlab_project_integration_jira`
- `gitlab_integration_mattermost` renamed to `gitlab_project_integration_mattermost`
- `gitlab_integration_microsoft_teams` renamed to `gitlab_project_integration_microsoft_teams`
- `gitlab_integration_pipelines_email` renamed to `gitlab_project_integration_pipelines_email`
- `gitlab_integration_redmine` renamed to `gitlab_project_integration_redmine`
- `gitlab_integration_telegram` renamed to `gitlab_project_integration_telegram`
- `gitlab_label` renamed to `gitlab_project_label`
- `gitlab_project_mirror` renamed to `gitlab_project_push_mirror`

## Resource replacements

- `gitlab_runner` replaced by `gitlab_user_runner`.
  - This switches to the [newer authentication method](https://docs.gitlab.com/runner/register/#register-with-a-runner-authentication-token) for registering runners.

## Attribute swaps

### Datasources

- `gitlab_projects._link` renamed to `gitlab_projects.links`.
  - Can be directly replaced with no other changes required.

### Resources

- `gitlab_integration_slack.notify_only_default_branch` switch to `gitlab_integration_slack.branches_to_be_notified`.
  - If `notify_only_default_branch` was `false`, set `branches_to_be_notified` to `all`.
  - If `notify_only_default_branch` was `true`, set `branches_to_be_notified` to `default`.
- `gitlab_project_share_group.access_level` switch to `gitlab_project_share_group.group_access`
  - Can be directly replaced with no other changes required.
- `gitlab_project`
  - For all of the following, replace `true` with `enabled` and `false` with `disabled`:
    - `issues_enabled` switch to `issues_access_level`
    - `merge_requests_enabled` switch to `merge_requests_access_level`
    - `pipelines_enabled` switch to `builds_access_level`
    - `wiki_enabled` switch to `wiki_access_level`
    - `snippets_enabled` switch to `snippets_access_level`
    - `container_registry_enabled` switch to `container_registry_access_level`
  - `restrict_user_defined_variables` switch to `ci_pipeline_variables_minimum_override_role`
    - If `restrict_user_defined_variables` was `false`, set `ci_pipeline_variables_minimum_override_role` to `developer`.
    - If `restrict_user_defined_variables` was `true`, set `ci_pipeline_variables_minimum_override_role` to `maintainer`.
  - `tags` switch to `topics`
    - Can be directly replaced with no other changes required.
  - `public_builds` switch to `public_jobs`
    - Can be directly replaced with no other changes required.
- `gitlab_application_settings.default_branch_protection` switch to `gitlab_application_settings.default_branch_protection_defaults`. As a rough guide:
  - If `default_branch_protection` was `0`:
  
    ```hcl
    default_branch_protection_defaults {
        allowed_to_push = [30] # Developer
        allowed_to_merge = [30] # Developer
        allow_force_push = true
    }
    ```

  - If `default_branch_protection` was `1`:
  
    ```hcl
    default_branch_protection_defaults {
        allowed_to_push = [30] # Developer
        allowed_to_merge = [40] # Maintainer
        allow_force_push = false
    }
    ```

  - If `default_branch_protection` was `2`:
  
    ```hcl
    default_branch_protection_defaults {
        allowed_to_push = [40] # Maintainer
        allowed_to_merge = [40] # Maintainer
        allow_force_push = false
    }
    ```

  - If `default_branch_protection` was `3`:
  
    ```hcl
    default_branch_protection_defaults {
        allowed_to_push = []
        allowed_to_merge = [40] # Maintainer
        allow_force_push = false
    }
    ```

## Attributes replaced by new resources

- `gitlab_project.approvals_before_merge` replaced by `gitlab_project_approval_rule`
  - All projects have a default approval rule, regardless of whether `approvals_before_merge` is in use.
  - By default, `gitlab_project_approval_rule` will automatically import the default approval rule.
  - Remove `gitlab_project.approvals_before_merge`.
  - Add the `gitlab_project_approval_rule` resource, and set attribute `approvals_required` to the value that was stored in `approvals_before_merge`.
  - Apply the changes
  - During the update process, the approvers total will be set to zero by the project resource for a short time. The approval rule resource will then import the rule and update the approvers total to the desired amount in the same apply operation.
  - Example old config:

    ```hcl
    resource "gitlab_project" "project" {
        approvals_before_merge = 2
    }
    ```

  - Example of new config:

    ```hcl
    resource "gitlab_project_approval_rule" "default_rule" {
        project            = gitlab_project.project.id
        name               = "Default"
        approvals_required = 2
    }
    ```

## Resources removed

These resources are for long deprecated features of GitLab.
The functionality will no longer be available in GitLab 19.0, so these resources will also be removed.

- `gitlab_group_cluster`
- `gitlab_instance_cluster`
- `gitlab_project_cluster`
