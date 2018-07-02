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
documentation](https://docs.gitlab.com/ce/ci/variables/README.html).


## Example Usage

```hcl
resource "gitlab_project_variable" "example" {
   project   = "example/project_with_variables"
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

## Attributes Reference

This resource does not currently export any attribute.
