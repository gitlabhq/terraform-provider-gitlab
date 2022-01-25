resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_jira" "jira" {
  project  = gitlab_project.awesome_project.id
  url      = "https://jira.example.com"
  username = "user"
  password = "mypass"
}
