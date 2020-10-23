# gitlab\_pipeline\_schedule

This resource allows you to create and manage pipeline schedules.
For further information on clusters, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/project/pipelines/schedules.html).

## Example Usage

```hcl
resource "gitlab_pipeline_schedule" "example" {
   project     = "12345"
   description = "Used to schedule builds"
   ref         = "master"
   cron        = "0 1 * * *"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the schedule to.

* `description` - (Required, string) The description of the pipeline schedule.

* `ref` - (Required, string) The branch/tag name to be triggered.

* `cron` - (Required, string) 	The cron (e.g. `0 1 * * *`).

* `cron_timezone` - (Optional, string) The timezone.

* `active` - (Optional, bool) The activation of pipeline schedule. If false is set, the pipeline schedule will deactivated initially.

## Importing pipeline schedules

You can import a pipeline schedule state using `terraform import <resource> <project_id>:<pipeline_schedule_id>`.

    terraform import gitlab_pipeline_schedule.example 123456789:13
