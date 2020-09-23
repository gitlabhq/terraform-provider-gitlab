# gitlab\_project\_variable

This resource allows you to create and manage CI/CD variables for your GitLab projects.
For further information on variables, consult the [gitlab
documentation](https://docs.gitlab.com/ce/ci/variables/README.html#variables).

~> **Important:** If your GitLab version is older than 13.4, you may see nondeterministic behavior
when updating or deleting `gitlab_project_variable` resources with non-unique keys, for example if
there is another variable with the same key and different environment scope. See
[this GitLab issue](https://gitlab.com/gitlab-org/gitlab/-/issues/9912).

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

* `variable_type` - (Optional, string)  The type of a variable. Available types are: env_var (default) and file.

* `protected` - (Optional, boolean) If set to `true`, the variable will be passed only to pipelines running on protected branches and tags. Defaults to `false`.

* `masked` - (Optional, boolean) If set to `true`, the variable will be masked if it would have been written to the logs. Defaults to `false`.

* `environment_scope` -  (Optional, string) The environment_scope of the variable. Defaults to `*`.

## Import

GitLab project variables can be imported using an id made up of `project:key:environment_scope`, e.g.

```
$ terraform import gitlab_project_variable.example '12345:project_variable_key:*'
```
