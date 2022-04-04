resource "gitlab_group" "this" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

resource "gitlab_project" "this" {
  name                   = "example"
  namespace_id           = gitlab_group.this.id
  initialize_with_readme = true
}

resource "gitlab_project_environment" "this" {
  project      = gitlab_project.this.id
  name         = "example"
  external_url = "www.example.com"
}
