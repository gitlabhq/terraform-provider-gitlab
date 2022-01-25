resource "gitlab_project" "awesome_project" {
  name             = "awesome_project"
  description      = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_pipelines_email" "email" {
  project                      = gitlab_project.awesome_project.id
  recipients                   = ["gitlab@user.create"]
  notify_only_broken_pipelines = true
  branches_to_be_notified      = "all"
}
