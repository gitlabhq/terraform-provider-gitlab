---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_group_service_account_access_token Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_group_service_account_access_token resource allows to manage the lifecycle of a group service account access token.
  ~> Use of the timestamp() function with expires_at will cause the resource to be re-created with every apply, it's recommended to use plantimestamp() or a static value instead.
  Upstream API: GitLab API docs https://docs.gitlab.com/ee/api/group_service_accounts.html#create-a-personal-access-token-for-a-service-account-user
---

# gitlab_group_service_account_access_token (Resource)

The `gitlab_group_service_account_access_token` resource allows to manage the lifecycle of a group service account access token.

~> Use of the `timestamp()` function with expires_at will cause the resource to be re-created with every apply, it's recommended to use `plantimestamp()` or a static value instead.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/group_service_accounts.html#create-a-personal-access-token-for-a-service-account-user)

## Example Usage

```terraform
resource "gitlab_group" "example" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

resource "gitlab_group_service_account" "example-sa" {
  group    = gitlab_group.example.id
  name     = "example-name"
  username = "example-username"
}

resource "gitlab_group_service_account_access_token" "example-sa-token" {
  group      = gitlab_group.example.id
  user_id    = gitlab_group_service_account.example-sa.id
  name       = "Example personal access token"
  expires_at = "2020-03-14"

  scopes = ["api"]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `group` (String) The ID or URL-encoded path of the group containing the service account. Must be a top level group.
- `name` (String) The name of the personal access token.
- `scopes` (Set of String) The scopes of the group service account access token. valid values are: `api`, `read_api`, `read_registry`, `write_registry`, `read_repository`, `write_repository`, `create_runner`, `manage_runner`, `ai_features`, `k8s_proxy`, `read_observability`, `write_observability`
- `user_id` (Number) The ID of a service account user.

### Optional

- `expires_at` (String) The personal access token expiry date. When left blank, the token follows the standard rule of expiry for personal access tokens.

### Read-Only

- `active` (Boolean) True if the token is active.
- `created_at` (String) Time the token has been created, RFC3339 format.
- `id` (String) The ID of the group service account access token.
- `revoked` (Boolean) True if the token is revoked.
- `token` (String, Sensitive) The token of the group service account access token. **Note**: the token is not available for imported resources.

## Import

Import is supported using the following syntax:

```shell
# You can import a service account access token using `terraform import <resource> <id>`.  The
# `id` is in the form of <group_id>:<service_account_id>:<access_token_id>
# Importing an access token does not import the access token value.
terraform import gitlab_group_service_account_access_token.example 1:2:3
```
