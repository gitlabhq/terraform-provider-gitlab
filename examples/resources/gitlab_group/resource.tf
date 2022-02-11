resource "gitlab_group" "example" {
  name        = "example"
  path        = "example"
  description = "An example group"
}

# Create a project in the example group
resource "gitlab_project" "example" {
  name         = "example"
  description  = "An example project"
  namespace_id = gitlab_group.example.id
}
