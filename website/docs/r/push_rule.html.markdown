---
layout: "gitlab"
page_title: "GitLab: gitlab_push_rule"
sidebar_current: "docs-gitlab-resource-push-rule"
description: |-
  Creates and manages push rules for GitLab projects
---

# gitlab\_push\_rule

This resource allows you to create and manage push rules for your GitLab projects.
For further information on variables, consult the [gitlab
documentation](https://docs.gitlab.com/ee/push_rules/push_rules.html).


## Example Usage

```hcl
resource "gitlab_push_rule" "example" {
   project   = "12345"
   max_file_size = 100
   commit_message_regex = "JIRA\-\d+"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the hook to.

* `commit_message_regex` - (Optional, string) All commit messages must match this regular expression to be pushed.

* `prevent_secrets` - (Optional, boolean) GitLab will reject any files that are likely to contain secrets. The list of file names we reject is available in the documentation.

* `max_file_size` - (Optional, integer) Pushes that contain added or updated files that exceed this file size are rejected. Set to 0 to allow files of any size.

* `deny_delete_tag` - (Optional, boolean) Tags can still be deleted through the web UI.

* `member_check` - (Optional, boolean) Restrict commits by author (email) to existing GitLab users.

* `branch_name_regex` - (Optional, string) All branch names must match this regular expression to be pushed. If this field is empty it allows any branch name.

* `author_email_regex` - (Optional, string) All commit author's email must match this regular expression to be pushed. If this field is empty it allows any email.

* `file_name_regex` - (Optional, string) All commited filenames must not match this regular expression to be pushed. If this field is empty it allows any filenames.

## Import

GitLab push rules can be imported using an id made up of `projectid`, e.g.

```
$ terraform import gitlab_push_rule.example 12345
```
