---
layout: "gitlab"
page_title: "GitLab: gitlab_datasource_project"
sidebar_current: "docs-gitlab-data-source-project"
description: |-
  View information about a project
---

# gitlab\_datasource_project

datasource_project provides details about a specific project in the gitlab provider. The results include the name of the project, path, description, default branch, etc.

## Example Usage

```hcl
data "gitlab_project" "test" {
	name = "Test Project"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project.

## Attributes Reference

The following attributes are exported:

* `path` - The path of the repository.

* `namespace_id` - The namespace (group or user) of the project. Defaults to your user.
  See [`gitlab_group`](group.html) for an example.

* `description` - A description of the project.

* `default_branch` - The default branch for the project.

* `issues_enabled` - Enable issue tracking for the project.

* `merge_requests_enabled` - Enable merge requests for the project.

* `wiki_enabled` - Enable wiki for the project.

* `snippets_enabled` - Enable snippets for the project.

* `visibility_level` -  Repositories are created as private by default.

* `id` - Integer that uniquely identifies the project within the gitlab install.

* `ssh_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via SSH.

* `http_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via HTTP.

* `web_url` - URL that can be used to find the project in a browser.



