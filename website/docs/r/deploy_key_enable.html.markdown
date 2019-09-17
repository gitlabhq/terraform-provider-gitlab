---
layout: "gitlab"
page_title: "GitLab: gitlab_deploy_key_enable"
sidebar_current: "docs-gitlab-resource-deploy_key_enable"
description: |-
  Enable a pre-existing deploy key in the project
---

# gitlab\_deploy\_key\_enable

This resource allows you to enable pre-existing deploy keys for your GitLab projects.

**the GITLAB KEY_ID for the deploy key must be known**

## Example Usage

```hcl
resource "gitlab_deploy_key_enable" "example" {
  project = "12345"
  key_id  = "67890"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the deploy key to.

* `key_id` - (Required, string) The Gitlab key id for the pre-existing deploy key

## Import

GitLab enabled deploy keys can be imported using an id made up of `{project_id}:{deploy_key_id}`, e.g.

```
$ terraform import gitlab_deploy_key_enable.example 12345:67890
```
