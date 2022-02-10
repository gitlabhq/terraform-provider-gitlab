resource "gitlab_branch" "this" {
  project = "12345"
  name    = "develop"
  ref     = "main"
}
