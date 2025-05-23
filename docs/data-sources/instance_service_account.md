---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_instance_service_account Data Source - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_instance_service_account data source retrieves information about a gitlab service account.
  ~> In order for a user to create a user account, they must have admin privileges at the instance level. This makes this feature unavailable on gitlab.com
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/user_service_accounts/
---

# gitlab_instance_service_account (Data Source)

The `gitlab_instance_service_account` data source retrieves information about a gitlab service account.

~> In order for a user to create a user account, they must have admin privileges at the instance level. This makes this feature unavailable on `gitlab.com`

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/user_service_accounts/)

## Example Usage

```terraform
data "gitlab_instance_service_account" "example" {
  service_account_id = "123"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `service_account_id` (String) The service account id.

### Read-Only

- `email` (String) The email of the user.
- `id` (String) The ID of this Terraform resource. This matches the service account id.
- `name` (String) The name of the user.
- `username` (String) The username of the user.
