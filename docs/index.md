# GitLab Provider

The GitLab provider is used to interact with GitLab group or user resources.

It needs to be configured with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources.

## Example Usage

```hcl
# Configure the GitLab Provider
provider "gitlab" {
    token = var.gitlab_token
}

# Add a project owned by the user
resource "gitlab_project" "sample_project" {
    name = "example"
}

# Add a hook to the project
resource "gitlab_project_hook" "sample_project_hook" {
    project = gitlab_project.sample_project.id
    url = "https://example.com/project_hook"
}

# Add a variable to the project
resource "gitlab_project_variable" "sample_project_variable" {
    project = gitlab_project.sample_project.id
    key = "project_variable_key"
    value = "project_variable_value"
}

# Add a deploy key to the project
resource "gitlab_deploy_key" "sample_deploy_key" {
    project = gitlab_project.sample_project.id
    title = "terraform example"
    key = "ssh-rsa AAAA..."
}

# Add a group
resource "gitlab_group" "sample_group" {
    name = "example"
    path = "example"
    description = "An example group"
}

# Add a project to the group - example/example
resource "gitlab_project" "sample_group_project" {
    name = "example"
    namespace_id = gitlab_group.sample_group.id
}
```

## Argument Reference

The following arguments are supported in the `provider` block:

* `token` - (Required) The OAuth2 token or project/personal access token used to connect to GitLab.
  It must be provided, but it can also be sourced from the `GITLAB_TOKEN` environment variable.

* `base_url` - (Optional) This is the target GitLab base API endpoint. Providing a value is a
  requirement when working with GitLab CE or GitLab Enterprise e.g. `https://my.gitlab.server/api/v4/`.
  It is optional to provide this value and it can also be sourced from the `GITLAB_BASE_URL` environment variable.
  The value must end with a slash.

* `cacert_file` - (Optional) This is a file containing the ca cert to verify the gitlab instance.  This is available
  for use when working with GitLab CE or Gitlab Enterprise with a locally-issued or self-signed certificate chain.

* `insecure` - (Optional; boolean, defaults to false) When set to true this disables SSL verification of the connection to the
  GitLab instance.

* `client_cert` - (Optional) File path to client certificate when GitLab instance is behind company proxy. File  must contain PEM encoded data.

* `client_key` - (Optional) File path to client key when GitLab instance is behind company proxy. File must contain PEM encoded data. Required when `client_cert` is set.
