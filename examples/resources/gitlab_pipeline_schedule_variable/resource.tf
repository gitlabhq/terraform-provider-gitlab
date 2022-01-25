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
