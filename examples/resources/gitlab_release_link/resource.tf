# Create a project
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = gitlab_group.example.id
}

# Create a tag
resource "gitlab_project_tag" "example" {
  name    = "example"
  ref     = gitlab_project.example.default_branch
  project = gitlab_project.example.id
}

resource "gitlab_release_link" "example" {
  project  = gitlab_project.example.id
  tag_name = gitlab_project_tag.example.name
  name     = "test"
  url      = "https://test/"
}
