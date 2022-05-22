# Create a project for the milestone to use
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = gitlab_group.example.id
}

resource "gitlab_project_milestone" "example" {
  project = gitlab_project.example.id
  title   = "example"
}
