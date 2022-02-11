resource "gitlab_project" "awesome_project" {
  name             = "awesome_project"
  description      = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_microsoft_teams" "teams" {
  project     = gitlab_project.awesome_project.id
  webhook     = "https://testurl.com/?token=XYZ"
  push_events = true
}
