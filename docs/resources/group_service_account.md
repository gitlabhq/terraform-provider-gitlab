---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_group_service_account Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_group_service_account resource allows creating a GitLab group service account.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/ee/api/group_service_accounts.html
---

# gitlab_group_service_account (Resource)

The `gitlab_group_service_account` resource allows creating a GitLab group service account.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/group_service_accounts.html)

## Example Usage

```terraform
# This must be a top-level group
resource "gitlab_group" "example" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

# The service account against the top-level group
resource "gitlab_group_service_account" "example_sa" {
  group    = gitlab_group.example.id
  name     = "example-name"
  username = "example-username"
}

# Group to assign the service account to. Can be the same top-level group resource as above, or a subgroup of that group.
resource "gitlab_group" "example_subgroup" {
  name        = "subgroup"
  path        = "example/subgroup"
  description = "An example subgroup"
}

# To assign the service account to a group
resource "gitlab_group_membership" "example_membership" {
  group_id     = gitlab_group.example_subgroup.id
  user_id      = gitlab_group_service_account.example_sa.service_account_id
  access_level = "developer"
  expires_at   = "2020-03-14"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group` (String) The ID or URL-encoded path of the group that the service account is created in. Must be a top level group.

### Optional

- `name` (String) The name of the user. If not specified, the default Service account user name is used.
- `username` (String) The username of the user. If not specified, it’s automatically generated.

### Read-Only

- `id` (String) The ID of this Terraform resource. In the format of `<group>:<service_account_id>`.
- `service_account_id` (String) The service account id.

## Import

Starting in Terraform v1.5.0 you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_group_service_account`. For example:
```terraform
import {
  to = gitlab_group_service_account.example
  id = "see CLI command below for ID"
}
```

Import using the CLI is supported using the following syntax:

```shell
# You can import a group service account using `terraform import <resource> <id>`.  The
# `id` is in the form of <group_id>:<service_account_id>
terraform import gitlab_group_service_account.example example
```