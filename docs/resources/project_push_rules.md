---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_project_push_rules Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_project_push_rules resource allows to manage the lifecycle of push rules on a project.
  ~> This resource will compete with the gitlab_project resource if push rules are also defined as
  part of that resource, since this resource will take over ownership of the project push rules created for the referenced project.
  It is recommended to define push rules using this resource OR in the gitlab_project resource,
  but not in both as it may result in terraform identifying changes with every "plan" operation.
  -> This resource requires a GitLab Enterprise instance with a Premium license to set the push rules on a project.
  Upstream API: GitLab API docs https://docs.gitlab.com/ee/api/projects.html#push-rules
---

# gitlab_project_push_rules (Resource)

The `gitlab_project_push_rules` resource allows to manage the lifecycle of push rules on a project.

~> This resource will compete with the `gitlab_project` resource if push rules are also defined as 
   part of that resource, since this resource will take over ownership of the project push rules created for the referenced project.
   It is recommended to define push rules using this resource OR in the `gitlab_project` resource, 
   but not in both as it may result in terraform identifying changes with every "plan" operation.

-> This resource requires a GitLab Enterprise instance with a Premium license to set the push rules on a project.

**Upstream API**: [GitLab API docs](https://docs.gitlab.com/ee/api/projects.html#push-rules)

## Example Usage

```terraform
resource "gitlab_project_push_rules" "sample" {
  project                       = 42
  author_email_regex            = "@gitlab.com$"
  branch_name_regex             = "(feat|fix)\\/*"
  commit_committer_check        = true
  commit_committer_name_check   = true
  commit_message_negative_regex = "ssh\\:\\/\\/"
  commit_message_regex          = "(feat|fix):.*"
  deny_delete_tag               = false
  file_name_regex               = "(jar|exe)$"
  max_file_size                 = 4
  member_check                  = true
  prevent_secrets               = true
  reject_unsigned_commits       = false
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project` (String) The ID or URL-encoded path of the project.

### Optional

- `author_email_regex` (String) All commit author emails must match this regex, e.g. `@my-company.com$`.
- `branch_name_regex` (String) All branch names must match this regex, e.g. `(feature|hotfix)\/*`.
- `commit_committer_check` (Boolean) Users can only push commits to this repository that were committed with one of their own verified emails.
- `commit_committer_name_check` (Boolean) Users can only push commits to this repository if the commit author name is consistent with their GitLab account name.
- `commit_message_negative_regex` (String) No commit message is allowed to match this regex, e.g. `ssh\:\/\/`.
- `commit_message_regex` (String) All commit messages must match this regex, e.g. `Fixed \d+\..*`.
- `deny_delete_tag` (Boolean) Deny deleting a tag.
- `file_name_regex` (String) All committed filenames must not match this regex, e.g. `(jar|exe)$`.
- `max_file_size` (Number) Maximum file size (MB).
- `member_check` (Boolean) Restrict commits by author (email) to existing GitLab users.
- `prevent_secrets` (Boolean) GitLab will reject any files that are likely to contain secrets.
- `reject_non_dco_commits` (Boolean) Reject commit when it’s not DCO certified.
- `reject_unsigned_commits` (Boolean) Reject commit when it’s not signed.

### Read-Only

- `id` (String) The ID of this Terraform resource.

## Import

Starting in Terraform v1.5.0 you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_project_push_rules`. For example:
```terraform
import {
  to = gitlab_project_push_rules.example
  id = "see CLI command below for ID"
}
```

Import using the CLI is supported using the following syntax:

```shell
# Gitlab project push rules can be imported with a key composed of `<project_id>`, e.g.
terraform import gitlab_project_push_rules.sample "42"
```