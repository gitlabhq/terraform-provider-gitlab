---
layout: "gitlab"
page_title: "GitLab: gitlab_group_variable"
sidebar_current: "docs-gitlab-resource-group-variable"
description: |-
  Creates and manages CI/CD variables for GitLab groups
---

# gitlab\_group\_variable

This resource allows you to create and manage CI/CD variables for your GitLab groups.
For further information on variables, consult the [gitlab
documentation](https://docs.gitlab.com/ce/ci/variables/README.html#variables).


## Example Usage

```hcl
resource "gitlab_group_variable" "example" {
   group     = "12345"
   key       = "group_variable_key"
   value     = "group_variable_value"
   protected = false
}
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, string) The name or id of the group to add the hook to.

* `key` - (Required, string) The name of the variable.

* `value` - (Required, string) The value of the variable.

* `variable_type` - (Optional, string)  The type of a variable. Available types are: env_var (default) and file.

* `protected` - (Optional, boolean) If set to `true`, the variable will be passed only to pipelines running on protected branches and tags. Defaults to `false`.

## Import

GitLab group variables can be imported using an id made up of `groupid:variablename`, e.g.

```
$ terraform import gitlab_group_variable.example 12345:group_variable_key
```
