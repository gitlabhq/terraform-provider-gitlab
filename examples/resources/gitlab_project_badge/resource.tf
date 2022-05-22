resource "gitlab_project" "foo" {
  name = "foo-project"
}

resource "gitlab_project_badge" "example" {
  project   = gitlab_project.foo.id
  link_url  = "https://example.com/badge-123"
  image_url = "https://example.com/badge-123.svg"
  name      = "badge-123"
}
