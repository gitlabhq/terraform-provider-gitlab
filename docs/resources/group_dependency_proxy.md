---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_group_dependency_proxy Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_group_dependency_proxy resource allows managing the group docker dependency proxy. More than one dependency proxy per group will conflict with each other.
  If you're looking to manage the project-level package dependency proxy, see the gitlab_project_package_registry_proxy resource instead.
  Upstream API: GitLab GraphQL API docs https://docs.gitlab.com/api/graphql/reference/#mutationupdatedependencyproxysettings
---

# gitlab_group_dependency_proxy (Resource)

The `gitlab_group_dependency_proxy` resource allows managing the group docker dependency proxy. More than one dependency proxy per group will conflict with each other.

If you're looking to manage the project-level package dependency proxy, see the `gitlab_project_package_registry_proxy` resource instead.

**Upstream API**: [GitLab GraphQL API docs](https://docs.gitlab.com/api/graphql/reference/#mutationupdatedependencyproxysettings)

## Example Usage

```terraform
resource "gitlab_group_dependency_proxy" "foo" {
  group = "1234"

  enabled  = true
  identity = "newidentity"
  secret   = "somesecret"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group` (String) The ID or URL-encoded path of the group.

### Optional

- `enabled` (Boolean) Indicates whether the proxy is enabled.
- `identity` (String) Identity credential used to authenticate with Docker Hub when pulling images. Can be a username (for password or personal access token (PAT)) or organization name (for organization access token (OAT)).
- `secret` (String, Sensitive) Secret credential used to authenticate with Docker Hub when pulling images. Can be a password, personal access token (PAT), or organization access token (OAT). Cannot be imported.

### Read-Only

- `id` (String) The ID of this Terraform resource. In the format of `<group-id>`.

## Import

Starting in Terraform v1.5.0, you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_group_dependency_proxy`. For example:

```terraform
import {
  to = gitlab_group_dependency_proxy.example
  id = "see CLI command below for ID"
}
```

Importing using the CLI is supported with the following syntax:

```shell
# You can import a group dependency proxy using the group id. e.g. `{group-id}`
# "secret" will not populate when importing the dependency proxy, but will still
# be required in the configuration.
terraform import gitlab_group_dependency_proxy.foo 42
```
