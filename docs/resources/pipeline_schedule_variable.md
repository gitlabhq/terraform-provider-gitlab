# gitlab\_pipeline\_schedule\_variable

This resource allows you to create and manage variables for pipeline schedules.

## Example Usage

```hcl
resource "gitlab_pipeline_schedule" "example" {
   project     = "12345"
   description = "Used to schedule builds"
   ref         = "master"
   cron        = "0 1 * * *"
}

resource "gitlab_pipeline_schedule_variable" "example" {
  project              = "${gitlab_pipeline_schedule.project}"
  pipeline_schedule_id = "${gitlab_pipeline_schedule.id}"
  key                  = "EXAMPLE_KEY"
  value                = "example"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The id of the project to add the schedule to.

* `pipeline_schedule_id` - (Required, string) The id of the pipeline schedule.

* `key` - (Required, string) Name of the variable.

* `value` - (Required, string) Value of the variable.

## Import

Pipeline schedule variables can be imported using an id made up of `project_id:pipeline_schedule_id:key`, e.g.

```
$ terraform import gitlab_pipeline_schedule_variable.example 123456789:13:mykey
```
