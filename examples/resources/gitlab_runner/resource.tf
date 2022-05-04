# Basic GitLab Runner
resource "gitlab_runner" "this" {
  token = "12345"
}

# GitLab Runner that runs only tagged jobs
resource "gitlab_runner" "tagged_only" {
  token       = "12345"
  description = "I only run tagged jobs"

  run_untagged = "false"
  tag_list     = ["tag_one", "tag_two"]
}

# GitLab Runner that only runs on protected branches
resource "gitlab_runner" "protected" {
  token       = "12345"
  description = "I only run protected jobs"

  access_level = "ref_protected"
}