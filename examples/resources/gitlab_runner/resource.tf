# Basic GitLab Group Runner
resource "gitlab_group" "my_group" {
  name        = "my runner"
  description = "group that holds the runners"
}
resource "gitlab_runner" "basic_runner" {
  registration_token = gitlab_group.my_group.runners_token
}

# GitLab Runner that runs only tagged jobs
resource "gitlab_runner" "tagged_only" {
  registration_token = gitlab_group.my_group.runners_token
  description        = "I only run tagged jobs"

  run_untagged = "false"
  tag_list     = ["tag_one", "tag_two"]
}

# GitLab Runner that only runs on protected branches
resource "gitlab_runner" "protected" {
  registration_token = gitlab_group.my_group.runners_token
  description        = "I only run protected jobs"

  access_level = "ref_protected"
}

# Generate a `config.toml` file that you can use to create a runner
# This is the typical workflow for this resource, using it to create an authentication_token which can then be used
# to generate the `config.toml` file to prevent re-registering the runner every time new hardware is created.

resource "gitlab_group" "my_custom_group" {
  name        = "my custom runner"
  description = "group that holds the custom runners"
}

resource "gitlab_runner" "my_runner" {
  registration_token = gitlab_group.my_custom_group.runners_token
}

# This creates a configuration for a local "shell" runner, but can be changed to generate whatever is needed.
# Place this configuration file on a server at `/etc/gitlab-runner/config.toml`, then run `gitlab-runner start`.
# See https://docs.gitlab.com/runner/configuration/advanced-configuration.html for more information.
resource "local_file" "config" {
  filename = "${path.module}/config.toml"
  content  = <<CONTENT
  concurrent = 1

  [[runners]]
    name = "Hello Terraform"
    url = "https://example.gitlab.com/"
    token = "${gitlab_runner.my_runner.authentication_token}"
    executor = "shell"
    
  CONTENT
}

