# gitlab\_project

Provide details about a specific project in the gitlab provider. The results include the name of the project, path, description, default branch, etc.

## Example Usage

```hcl
data "gitlab_project" "example" {
  id = 30
}
```

```hcl
data "gitlab_project" "example" {
  id = "foo/bar/baz"
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The integer or path with namespace that uniquely identifies the project within the gitlab install.

## Attributes Reference

The following attributes are exported:

* `path` - The path of the repository.

* `path_with_namespace` - The path of the repository with namespace.

* `namespace_id` - The namespace (group or user) of the project. Defaults to your user.
  See [`gitlab_group`](../resources/group) for an example.

* `description` - A description of the project.

* `default_branch` - The default branch for the project.

* `request_access_enabled` - Allow users to request member access.

* `issues_enabled` - Enable issue tracking for the project.

* `merge_requests_enabled` - Enable merge requests for the project.

* `pipelines_enabled` - Enable pipelines for the project.

* `wiki_enabled` - Enable wiki for the project.

* `snippets_enabled` - Enable snippets for the project.

* `lfs_enabled` - Enable LFS for the project.

* `visibility_level` -  Repositories are created as private by default.

* `id` - Integer that uniquely identifies the project within the gitlab install.

* `ssh_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via SSH.

* `http_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via HTTP.

* `web_url` - URL that can be used to find the project in a browser.

* `runners_token` - Registration token to use during runner setup.

* `archived` - Whether the project is in read-only mode (archived).

* `remove_source_branch_after_merge` - Enable `Delete source branch` option by default for all new merge requests

* `packages_enabled` - Enable packages repository for the project.

* `push_rules` Push rules for the project (documented below).

## Nested Blocks

### push_rules

For information on push rules, consult the [GitLab documentation](https://docs.gitlab.com/ce/push_rules/push_rules.html#push-rules).

#### Attributes

* `author_email_regex` - All commit author emails must match this regex, e.g. `@my-company.com$`.

* `branch_name_regex` - All branch names must match this regex, e.g. `(feature|hotfix)\/*`.

* `commit_message_regex` - All commit messages must match this regex, e.g. `Fixed \d+\..*`.

* `commit_message_negative_regex` - No commit message is allowed to match this regex, for example `ssh\:\/\/`.

* `file_name_regex` - All commited filenames must not match this regex, e.g. `(jar|exe)$`.

* `commit_committer_check` - Users can only push commits to this repository that were committed with one of their own verified emails.

* `deny_delete_tag` - Deny deleting a tag.

* `member_check` - Restrict commits by author (email) to existing GitLab users.

* `prevent_secrets` - GitLab will reject any files that are likely to contain secrets.

* `reject_unsigned_commits` - Reject commit when it’s not signed through GPG.

* `max_file_size` - Maximum file size (MB).
