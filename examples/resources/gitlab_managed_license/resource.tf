resource "gitlab_project" "foo" {
  name             = "example project"
  description      = "Lorem Ipsum"
  visibility_level = "public"
}

resource "gitlab_managed_license" "mit" {
  project         = gitlab_project.foo.id
  name            = "MIT license"
  approval_status = "approved"
}
