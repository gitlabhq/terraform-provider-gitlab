# Create a project for the tag to use
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = gitlab_group.example.id
}

resource "gitlab_project_tag" "example" {
  name    = "example"
  ref     = "main"
  project = gitlab_project.example.id
}

