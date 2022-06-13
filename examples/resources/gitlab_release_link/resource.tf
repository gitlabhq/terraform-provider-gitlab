# Create a project
resource "gitlab_project" "example" {
  name        = "example"
  description = "An example project"
}

# Can create release link only to a tag associated with a release
resource "gitlab_release_link" "example" {
  project  = gitlab_project.example.id
  tag_name = "tag_name_associated_with_release"
  name     = "test"
  url      = "https://test/"
}
