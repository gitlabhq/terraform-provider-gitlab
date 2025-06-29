---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_pipeline_schedule Resource - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_pipeline_schedule resource allows to manage the lifecycle of a scheduled pipeline.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/pipeline_schedules/
---

# gitlab_pipeline_schedule (Resource)

The `gitlab_pipeline_schedule` resource allows to manage the lifecycle of a scheduled pipeline.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/pipeline_schedules/)

## Example Usage

```terraform
resource "gitlab_pipeline_schedule" "example" {
  project     = "12345"
  description = "Used to schedule builds"
  ref         = "refs/heads/main"
  cron        = "0 1 * * *"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `cron` (String) The cron (e.g. `0 1 * * *`).
- `description` (String) The description of the pipeline schedule.
- `project` (String) The name or id of the project to add the schedule to.
- `ref` (String) The branch/tag name to be triggered. This must be the full branch reference, for example: `refs/heads/main`, not `main`.

### Optional

- `active` (Boolean) The activation of pipeline schedule. If false is set, the pipeline schedule will deactivated initially.
- `cron_timezone` (String) The timezone.
- `take_ownership` (Boolean) When set to `true`, the user represented by the token running Terraform will take ownership of the scheduled pipeline prior to editing it. This can help when managing scheduled pipeline drift when other users are making changes outside Terraform.

### Read-Only

- `id` (String) The ID of this Terraform resource. In the format of `<project-id>:<pipeline-schedule-id>`.
- `owner` (Number) The ID of the user that owns the pipeline schedule.
- `pipeline_schedule_id` (Number) The pipeline schedule id.

## Import

Starting in Terraform v1.5.0, you can use an [import block](https://developer.hashicorp.com/terraform/language/import) to import `gitlab_pipeline_schedule`. For example:

```terraform
import {
  to = gitlab_pipeline_schedule.example
  id = "see CLI command below for ID"
}
```

Importing using the CLI is supported with the following syntax:

```shell
# GitLab pipeline schedules can be imported using an id made up of `{project_id}:{pipeline_schedule_id}`, e.g.
terraform import gitlab_pipeline_schedule.test 1:3
```
