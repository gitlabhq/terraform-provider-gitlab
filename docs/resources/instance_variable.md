# gitlab\_instance\_variable

This resource allows you to create and manage CI/CD variables for your GitLab instance.
For further information on variables, consult the [gitlab
documentation](https://docs.gitlab.com/ee/api/instance_level_ci_variables.html).

## Example Usage

```hcl
resource "gitlab_instance_variable" "example" {
   key       = "instance_variable_key"
   value     = "instance_variable_value"
   protected = false
   masked    = false
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required, string) The name of the variable.

* `value` - (Required, string) The value of the variable.

* `variable_type` - (Optional, string)  The type of a variable. Available types are: env_var (default) and file.

* `protected` - (Optional, boolean) If set to `true`, the variable will be passed only to pipelines running on protected branches and tags. Defaults to `false`.

* `masked` - (Optional, boolean) If set to `true`, the value of the variable will be hidden in job logs. The value must meet the [masking requirements](https://docs.gitlab.com/ee/ci/variables/#masked-variable-requirements). Defaults to `false`.

## Import

GitLab instance variables can be imported using an id made up of `variablename`, e.g.

```console
$ terraform import gitlab_instance_variable.example instance_variable_key
```
