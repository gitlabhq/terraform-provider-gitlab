resource "gitlab_project_freeze_period" "schedule" {
  project_id    = gitlab_project.foo.id
  freeze_start  = "0 23 * * 5"
  freeze_end    = "0 7 * * 1"
  cron_timezone = "UTC"
}
