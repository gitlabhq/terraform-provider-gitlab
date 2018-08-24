---
layout: "gitlab"
page_title: "GitLab: gitlab_project_variable"
sidebar_current: "docs-gitlab-resource-project-variable"
description: |-
  Creates and manages CI/CD variables for GitLab projects
---

# gitlab\_project\_variable

This resource allows you to create and manage CI/CD variables for your GitLab projects.
For further information on variables, consult the [gitlab
documentation](https://docs.gitlab.com/ce/ci/variables/README.html#variables).


## Example Usage

```hcl
resource "gitlab_project_variable" "example" {
   project   = "12345"
   key       = "project_variable_key"
   value     = "project_variable_value"
   protected = false
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the hook to.

* `key` - (Required, string) The name of the variable.

* `value` - (Required, string) The value of the variable.

* `protected` - (Optional, boolean) If set to `true`, the variable will be passed only to pipelines running on protected branches and tags. Defaults to `false`.

## Import

GitLab project variables can be imported using an id made up of `projectid:variablename`, e.g.

```
$ terraform import gitlab_group_membership.test 12345:project_variable_key
```