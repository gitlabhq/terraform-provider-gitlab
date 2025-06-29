---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_pipeline_schedule_variable Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_pipeline_schedule_variable resource allows to manage the lifecycle of a variable for a pipeline schedule.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/pipeline_schedules/#pipeline-schedule-variables
---

# gitlab_pipeline_schedule_variable (Resource)

The `gitlab_pipeline_schedule_variable` resource allows to manage the lifecycle of a variable for a pipeline schedule.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/pipeline_schedules/#pipeline-schedule-variables)

## Example Usage

```terraform
resource "gitlab_pipeline_schedule" "example" {
  project     = "12345"
  description = "Used to schedule builds"
  ref         = "master"
  cron        = "0 1 * * *"
}

resource "gitlab_pipeline_schedule_variable" "example" {
  project              = gitlab_pipeline_schedule.example.project
  pipeline_schedule_id = gitlab_pipeline_schedule.example.pipeline_schedule_id
  key                  = "EXAMPLE_KEY"
  value                = "example"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `key` (String) Name of the variable.
- `pipeline_schedule_id` (Number) The id of the pipeline schedule.
- `project` (String) The id of the project to add the schedule to.
- `value` (String) Value of the variable.

### Optional

- `variable_type` (String) The type of a variable. Available types are: `env_var`, `file`. Default is `env_var`.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Starting in Terraform v1.5.0, you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_pipeline_schedule_variable`. For example:

```terraform
import {
  to = gitlab_pipeline_schedule_variable.example
  id = "see CLI command below for ID"
}
```

Importing using the CLI is supported with the following syntax:

```shell
# Pipeline schedule variables can be imported using an id made up of `project_id:pipeline_schedule_id:key`, e.g.
terraform import gitlab_pipeline_schedule_variable.example 123456789:13:mykey
```
