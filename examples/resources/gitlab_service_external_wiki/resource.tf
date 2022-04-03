resource "gitlab_project" "awesome_project" {
  name             = "awesome_project"
  description      = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_external_wiki" "wiki" {
  project           = gitlab_project.awesome_project.id
  external_wiki_url = "https://MyAwesomeExternalWikiURL.com"
}
