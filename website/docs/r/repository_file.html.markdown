---
layout: "gitlab"
page_title: "Gitlab: gitlab_repository_file"
sidebar_current: "docs-gitlab-repository-file"
---

# gitlab\_repository\_file

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

resource "gitlab_repository_file" "this" {
    project        = gitlab_project.this.id
    file           = "meow.txt"
    content        = base64encode("Meow goes the cat")
    branch         = "master"
    author_name    = "Terraform"
    author_email   = "terraform@example.com"
    commit_message = "feature: add meow file"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The ID of the project.

* `file` - (Required) The full path of the file.

* `content` - (Required) Base64 encoded string.

* `branch` - (Required) Name of the branch to which to commit to.

* `author_name` - (Optional) Name of the commit author.

* `author_email` - (Optional) Email of the commit author.

* `commit_message` - (Required) Commit message.

## Attribute Reference

The resource exports the following attributes:

* `id` - The unique ID assigned to the file.
