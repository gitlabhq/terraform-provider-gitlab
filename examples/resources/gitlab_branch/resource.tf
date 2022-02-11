resource "gitlab_branch" "this" {
  project = "12345"
  branch  = "develop"
  ref     = "main"
}
