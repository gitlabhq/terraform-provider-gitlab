# gitlab\_repository\_file

This resource allows you to create and manage GitLab repository files.

**Limitations**:

The [GitLab Repository Files API](https://docs.gitlab.com/ee/api/repository_files.html)
can only create, update or delete a single file at the time.
The API will also
[fail with a `400`](https://docs.gitlab.com/ee/api/repository_files.html#update-existing-file-in-repository)
response status code if the underlying repository is changed while the API tries to make changes.
Therefore, it's recommended to make sure that you execute it with
[`-parallelism=1`](https://www.terraform.io/docs/cli/commands/apply.html#parallelism-n)
and that no other entity than the terraform at hand makes changes to the
underlying repository while it's executing.

## Example Usage

```hcl
resource "gitlab_group" "this" {
    name        = "example"
    path        = "example"
    description = "An example group"
}
resource "gitlab_project" "this" {
    name                   = "example"
    namespace_id           = gitlab_group.this.id
    initialize_with_readme = true
}
resource "gitlab_repository_file" "this" {
    project        = gitlab_project.this.id
    file_path      = "meow.txt"
    branch         = "main"
    content        = base64encode("Meow goes the cat")
    author_email   = "terraform@example.com"
    author_name    = "Terraform"
    commit_message = "feature: add meow file"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The ID of the project.

* `file_path` - (Required) The full path of the file.
                It must be relative to the root of the project without a leading slash `/`.

* `branch` - (Required) Name of the branch to which to commit to.

* `content` - (Required) base64 encoded file content.
              No other encoding is currently supported,
              because of a [GitLab API bug](https://gitlab.com/gitlab-org/gitlab/-/issues/342430).

* `commit_message` - (Required) Commit message.

* `start_branch` - (Optional) Name of the branch to start the new commit from.

* `author_email` - (Optional) Email of the commit author.

* `author_name` - (Optional) Name of the commit author.

## Attribute Reference

The resource exports the following attributes:

* `id` - The unique ID assigned to the file.

## Import

A Repository File can be imported using the following id pattern, e.g.

```
$ terraform import gitlab_repository_file.this <project-id>:<branch-name>:<file-path>
```
