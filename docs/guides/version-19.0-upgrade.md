---
page_title: "Terraform GitLab Provider Version 19.0 Upgrade Guide"
subcategory: "Upgrade Guides"
---

# Upgrade to Terraform GitLab Provider Version 19.0

The GitLab 19.0 major milestone introduced some breaking changes that this release addresses.
The provider also has some breaking changes which may require actions on the users side.
These are described below:

## Resource renames

These resources are identical but have been renamed.
You can migrate between the two using a state move operation.

Rename the resource reference in your code. For example:

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

Perform a state move using one of the available options:

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

~> To use the moved block across resources types like this, you will need at least [Terraform version 1.8](https://github.com/hashicorp/terraform/blob/v1.8.0/CHANGELOG.md) and GitLab Provider version 18.9.

### Renamed resources

- `gitlab_deploy_token` renamed to `gitlab_project_deploy_token`
- `gitlab_integration_custom_issue_tracker` renamed to `gitlab_project_integration_custom_issue_tracker`
- `gitlab_integration_emails_on_push` renamed to `gitlab_project_integration_emails_on_push`
- `gitlab_integration_external_wiki` renamed to `gitlab_project_integration_external_wiki`
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

## Resource gitlab_runner replacement

Replaced by `gitlab_user_runner`.
This switches to the [newer authentication method](https://docs.gitlab.com/runner/register/#register-with-a-runner-authentication-token) for registering runners.

## Attribute swaps

### Datasource gitlab_projects

The `gitlab_projects._link` attribute has been renamed to `gitlab_projects.links`.
This can be directly replaced with no other changes required.

### Resource gitlab_integration_slack

The `gitlab_integration_slack.notify_only_default_branch` attribute should be replaced with `gitlab_integration_slack.branches_to_be_notified`.

- If `notify_only_default_branch` was `false`, set `branches_to_be_notified` to `all`.
- If `notify_only_default_branch` was `true`, set `branches_to_be_notified` to `default`.

### Resource gitlab_project_share_group

The `gitlab_project_share_group.access_level` attribute should be replaced with `gitlab_project_share_group.group_access`.
This can be directly replaced with no other changes required.

### Resource gitlab_project

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

### Datasource gitlab_project

- For all of the following, `true` becomes `enabled` and `false` becomes `disabled`:
  - `issues_enabled` switch to `issues_access_level`
  - `merge_requests_enabled` switch to `merge_requests_access_level`
  - `pipelines_enabled` switch to `builds_access_level`
  - `wiki_enabled` switch to `wiki_access_level`
  - `snippets_enabled` switch to `snippets_access_level`
- `restrict_user_defined_variables` switch to `ci_pipeline_variables_minimum_override_role`
  - If `restrict_user_defined_variables` was `false`, `ci_pipeline_variables_minimum_override_role` becomes `developer`.
  - If `restrict_user_defined_variables` was `true`, `ci_pipeline_variables_minimum_override_role` becomes `maintainer`.
- `public_builds` switch to `public_jobs`
  - Can be directly replaced with no other changes required.

### Datasource gitlab_project_membership

If you were previously using `project_id` or `full_path`, they should be directly replaced with `project`.

### Datasource gitlab_repository_tree

If you were previously using `tree.id`, replace with `tree.node_id`.

### Resource gitlab_project_protected_environment

The `deploy_access_levels` attribute should be replaced by the `deploy_access_levels_attribute` attribute.

For example:

```hcl
  resource "gitlab_project_protected_environment" "this" {
    project     = 123
    environment = "production"

    deploy_access_levels {
      group_id = 456
    }

    deploy_access_levels {
      access_level = "developer"
    }
  }
```

Becomes:

```hcl
  resource "gitlab_project_protected_environment" "this" {
    project     = 123
    environment = "production"

    deploy_access_levels_attribute = [
      {
        group_id = 456
      },
      {
        access_level = "developer"
      }
  }
```

### Resource gitlab_application_settings

The `gitlab_application_settings.default_branch_protection` attribute should be replaced with `gitlab_application_settings.default_branch_protection_defaults`.
As a rough guide:

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

### Datasource gitlab_project_protected_branch

The `gitlab_project_protected_branch.merge_access_levels` and `gitlab_project_protected_branch.push_access_levels`
attributes have been migrated from a Block Set to an Attributes List and are now read-only.

### Datasource gitlab_project_protected_branches

The `gitlab_project_protected_branches.protected_branches` attribute has been migrated from a Block List
to an Attributes List and is read-only.

The `gitlab_project_protected_branch.protected_branches.merge_access_levels` and
`gitlab_project_protected_branch.protected_branches.push_access_levels` attributes
have been migrated from a Block Set to an Attributes List and are now read-only.

## Resource gitlab_branch_protection overhaul

The `gitlab_branch_protection` resource has been reworked to have separate attributes for CE and EE licenses.  This will cause existing
configurations to break.

CE users will continue to use the `gitlab_branch_protection.push_access_level` and `gitlab_branch_protection.merge_access_level` attributes.

For EE users, the `gitlab_branch_protection.push_access_level`, `gitlab_branch_protection.merge_access_level`, and
`gitlab_branch_protection.unprotect_access_level` attributes have been removed.  Their usage should be replaced with
`gitlab_branch_protection.allowed_to_push`, `gitlab_branch_protection.allowed_to_merge`, and
`gitlab_branch_protection.allowed_to_unprotect` attributes.

Example old EE config:

```hcl
resource "gitlab_branch_protection" "branchA" {
  project                      = "12345"
  branch                       = "branchA"
  push_access_level            = "developer"
  merge_access_level           = "developer"
  unprotect_access_level       = "developer"
  allow_force_push             = true
  code_owner_approval_required = true

  allowed_to_push {
    user_id = 5
  }
  allowed_to_merge {
    user_id = 37
  }
  allowed_to_unprotect {
    group_id = 42
  }
}
```

Example new EE config:

```hcl
resource "gitlab_branch_protection" "branchA" {
  project                      = "12345"
  branch                       = "branchA"
  allow_force_push             = true
  code_owner_approval_required = true

  allowed_to_push = [
    {
      user_id = 5
    },
    {
      access_level = "developer"
    }
  ]

  allowed_to_merge = [
    {
      user_id = 37
    },
    {
      access_level = "developer"
    }
  ]

  allowed_to_unprotect = [
    {
      group_id = 42
    },
    {
      access_level = "developer"
    }
  ]
}
```

## Resource gitlab_project.approvals_before_merge Replacement

The `gitlab_project.approvals_before_merge` attribute should be replaced with the `gitlab_project_approval_rule` resource.

Example old config:

```hcl
resource "gitlab_project" "project" {
    approvals_before_merge = 2
}
```

Example new config:

```hcl
resource "gitlab_project_approval_rule" "default_rule" {
    project            = gitlab_project.project.id
    name               = "Default"
    approvals_required = 2
}
```

### Approval Rule Migration Process

All projects have a default approval rule, regardless of whether `approvals_before_merge` is in use.
By default, `gitlab_project_approval_rule` will automatically import the default approval rule.

- Remove `gitlab_project.approvals_before_merge`.
- Add the `gitlab_project_approval_rule` resource, and set attribute `approvals_required` to the value that was stored in `approvals_before_merge`.
- Apply the changes

During the update process, the approvers total will be set to zero by the project resource for a short time.
The approval rule resource will then import the rule and update the approvers total to the desired amount in the same apply operation.

## Resource gitlab_project.mirror Replacement

The mirror attributes on `gitlab_project` should be replaced by equivalent attributes on the new `gitlab_project_pull_mirror` resource.

Example old config:

```hcl
resource "gitlab_project" "import_private_with_mirror" {
  name                                = "import-from-public-project"
  import_url                          = "https://gitlab.example.com/repo.git"
  import_url_username                 = "user"
  import_url_password                 = "pass"
  mirror                              = true
  mirror_trigger_builds               = true
  only_mirror_protected_branches      = true
  mirror_overwrites_diverged_branches = true
}
```

Example new config:

```hcl
resource "gitlab_project" "import_private_with_mirror" {
  name = "import-from-public-project"
}

resource "gitlab_project_pull_mirror" "mirror" {
  project                             = gitlab_project.import_private_with_mirror.id
  url                                 = "https://gitlab.example.com/repo.git"
  auth_user                           = "user"
  auth_password                       = "pass"
  mirror_trigger_builds               = true
  only_mirror_protected_branches      = true
  mirror_overwrites_diverged_branches = true
}
```

### Mirror Migration Process

- Remove any usage of the following attributes on the `gitlab_project` resource:
  - `import_url`
  - `import_url_username`
  - `import_url_password`
  - `mirror`
  - `mirror_trigger_builds`
  - `only_mirror_protected_branches`
  - `mirror_overwrites_diverged_branches`
- Add code for the new `gitlab_project_pull_mirror` resource with the following attribute values if set on the old resource:
  - `url` from `gitlab_project.import_url`
  - `auth_user` from `gitlab_project.import_url_username`
  - `auth_password` from `gitlab_project.import_url_password`
  - `mirror_trigger_builds` from `gitlab_project.mirror_trigger_builds`
  - `only_mirror_protected_branches` from `gitlab_project.only_mirror_protected_branches`
  - `mirror_overwrites_diverged_branches` from `gitlab_project.mirror_overwrites_diverged_branches`
- Run an apply.
  - As there is a relationship between the two resources, the `gitlab_project` resource will apply first. This will temporarily remove the mirror.
  - Then the new pull mirror resource will apply and add the mirror back in.

## Any problems upgrading?

Please [fill out an issue](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/work_items) in this project's issue tracker and someone from the community will respond as soon as they are available to help you.
