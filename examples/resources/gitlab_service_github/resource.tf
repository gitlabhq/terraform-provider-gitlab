resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_github" "github" {
  project        = gitlab_project.awesome_project.id
  token          = "REDACTED"
  repository_url = "https://github.com/gitlabhq/terraform-provider-gitlab"
}
