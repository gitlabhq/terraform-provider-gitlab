---
page_title: "Terraform GitLab Provider Version 15.7 Upgrade Guide"
---

# Upgrade to Terraform GitLab Provider Version 15.7

Because the GitLab Provider has moved to [GitLab.com](https://gitlab.com/gitlab-org/terraform-provider-gitlab),
the release cadence and versioning has been
[aligned](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/issues/1331) with the GitLab
[monthly self-managed release cadence](https://about.gitlab.com/handbook/engineering/releases/)
starting with 15.7 (22nd Dec 2022).

The version bump from [`3.20.0`](https://registry.terraform.io/providers/gitlabhq/gitlab/3.20.0) to
[`v15.7.0`](https://registry.terraform.io/providers/gitlabhq/gitlab/15.7.0) introduced a few breaking changes,
which are described below.


## Terraform version 1.0

The GitLab Provider upgraded to
[Terraform Protocol v6](https://developer.hashicorp.com/terraform/plugin/how-terraform-works#protocol-version-6),
which **requires at least Terraform 1.0**.

## Provider token is now sensitive

The `token` Provider argument is now marked as
[`sensitive`](https://developer.hashicorp.com/terraform/tutorials/configuration-language/sensitive-variables).
This affects current Provider configurations, which read the `token` value from another non-sensitive Terraform value,
like an output attribute or variable.

Because the `variable "gitlab_token"` declaration doesn't mark the variable as `sensitive`,
code like the following will break:

```hcl
variable "gitlab_token" {
  type = string
}

provider "gitlab" {
  token = var.gitlab_token
}
```
There are two ways to resolve this issue:

- If you have control over the token value, mark it as
[sensitive](https://developer.hashicorp.com/terraform/language/values/variables#suppressing-values-in-cli-output):

  ```hcl
  variable "gitlab_token" {
    type      = string
    sensitive = true
  }

  provider "gitlab" {
    token = var.gitlab_token
  }
  ```

- If you don't have control over the token value, use the
[`sensitive()`](https://developer.hashicorp.com/terraform/language/functions/sensitive) function to create
a *sensitive* copy of the value to use:

  ```hcl
  provider "gitlab" {
    token = sensitive(some_other_module.gitlab_token)
  }
  ```
