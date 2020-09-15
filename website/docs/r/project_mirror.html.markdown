---
layout: "gitlab"
page_title: "GitLab: gitlab_project_mirror"
sidebar_current: "docs-gitlab-resource-project_mirror"
description: |-
  Adds a target to a projects remote mirrors.
---

# gitlab\_project_mirror

This resource allows you to add a mirror target for the repository, all changes will be synced to the remote target.

## Example Usage

```hcl
resource "gitlab_project_mirror" "foo" {
  project = "1"
  url = "https://username:password@github.com/org/repository.git"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The id of the project.

* `url` - (Required) 	The URL of the remote repository to be mirrored.

* `enabled` - 	Determines if the mirror is enabled.

* `only_protected_branches` - Determines if only protected branches are mirrored.

* `keep_divergent_refs` - Determines if divergent refs are skipped.

## Import

GitLab project mirror can be imported using an id made up of `project_id:mirror_id`, e.g.


```
$ terraform import gitlab_project_mirror.foo "12345:1337"
```
