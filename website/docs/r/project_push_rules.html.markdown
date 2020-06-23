---
layout: "gitlab"
page_title: "GitLab: gitlab_project_push_rules"
sidebar_current: "docs-gitlab-resource-project-push-rules"
description: |-
  Creates and manages push rules for GitLab projects
---

# gitlab\_project\_push\_rules

This resource allows you to create and manage push rules for your GitLab projects.
For further information on push rules, consult the [gitlab
documentation](https://docs.gitlab.com/ce/push_rules/push_rules.html#push-rules).

## Example Usage

```hcl
resource "gitlab_project_push_rules" "example" {
  commit_message_regex = "^(feat|feature|fix|chore|docs|BREAKING_CHANGE):.*"
  prevent_secrets = true
  branch_name_regex = "^PROJ-\d+-.*"
  author_email_regex = "@my-company.com$"
  commit_committer_check = true
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the push rules to.

* `commit_message_regex` - (Optional, string) All commit messages must match this regex, e.g. "Fixed \d+\..*"

* `deny_delete_tag` - (Optional, bool) Deny deleting a tag

* `member_check` - (Optional, bool) Restrict commits by author (email) to existing GitLab users

* `prevent_secrets` - (Optional, bool) GitLab will reject any files that are likely to contain secrets

* `branch_name_regex` - (Optional, string) All branch names must match this regex, e.g. "(feature|hotfix)\/*"

* `author_email_regex` - (Optional, string) All commit author emails must match this regex, e.g. "@my-company.com$"

* `file_name_regex` - (Optional, string) All commited filenames must not match this regex, e.g. "(jar|exe)$"

* `max_file_size` - (Optional, int) Maximum file size (MB)

* `commit_committer_check` - (Optional, bool) Users can only push commits to this repository that were committed with one of their own verified emails

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the push rules by the GitLab server.

## Import

GitLab push rules can be imported using the project id or name, as when importing a project, e.g.

```
$ terraform import gitlab_project_push_rules.test richardc/example
```
