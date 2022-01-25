resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_slack" "slack" {
  project                    = gitlab_project.awesome_project.id
  webhook                    = "https://webhook.com"
  username                   = "myuser"
  push_events                = true
  push_channel               = "push_chan"
}
