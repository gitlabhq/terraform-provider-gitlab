# gitlab\_repository\_files

This resource allows you to create and manage GitLab repository files.

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

resource "gitlab_repository_files" "this" {
  project        = gitlab_project.this.id
  branch         = "main"
  author_email   = "terraform@example.com"
  author_name    = "Terraform"
  commit_message = "feature: add many files"

  dynamic "file" {
    for_each = [
      { file_path : "my-file-1.txt", content : "some content 1" },
      { file_path : "my-file-2.txt", content : "some content 2" },
      { file_path : "my-file-3.txt", content : "some content 3" },
    ]

    content {
      file_path = file.value.file_path
      content   = file.value.content
    }
  }
}

```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The ID of the project.

* `branch` - (Required) Name of the branch to which to commit to.

* `file` - (Required) The file object structure is described below.

* `commit_message` - (Required) Commit message.

* `start_branch` - (Optional) Name of the branch to start the new commit from.

* `author_email` - (Optional) Email of the commit author.

* `author_name` - (Optional) Name of the commit author.

The `file` block supports:

* `file_path` - (Required) The full path of the file. It must be relative to the root of the project without a leading
  slash `/`.

* `content` - (Required) The file content in clear text.

## Import

Repository files can be imported using the following id pattern where each file path is comma separated, e.g.

 ```
 $ terraform import gitlab_repository_files.this <project-id>:<branch-name>:<file-path1,file-path2,file-path3>
 ```