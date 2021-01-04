# gitlab_project_freeze_period

This resource allows you to create and manage freeze periods. For further information on freeze periods, consult the [gitlab documentation](https://docs.gitlab.com/ee/api/freeze_periods.html#create-a-freeze-period).

## Example Usage

```hcl
resource "gitlab_project_freeze_period" "schedule" {
	project = gitlab_project.foo.id
	freeze_start = "0 23 * * 5"
	freeze_end =  "0 7 * * 1"
	cron_timezone = "UTC"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The name or id of the project to add the schedule to.

* `freeze_start` - (Required,string) Start of the Freeze Period in cron format (e.g. `0 1 * * *`).

* `freeze_end` - (Required, string) End of the Freeze Period in cron format (e.g. `0 2 * * *`).

* `cron_timezone` - (Optional, string) The timezone.
