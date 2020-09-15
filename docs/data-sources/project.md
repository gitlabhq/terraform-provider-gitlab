---
layout: "gitlab"
page_title: "GitLab: gitlab_project"
sidebar_current: "docs-gitlab-data-source-project"
description: |-
  View information about a project
---

# gitlab\_project

Provides details about a specific project in the gitlab provider. The results include the name of the project, path, description, default branch, etc.

## Example Usage

```hcl
data "gitlab_project" "example" {
	id = 30
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The integer that uniquely identifies the project within the gitlab install.

## Attributes Reference

The following attributes are exported:

* `path` - The path of the repository.

* `path_with_namespace` - The path of the repository with namespace.

* `namespace_id` - The namespace (group or user) of the project. Defaults to your user.
  See [`gitlab_group`](../r/group.html) for an example.

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