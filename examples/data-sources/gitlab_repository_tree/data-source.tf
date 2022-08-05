data "gitlab_repository_tree" "this" {
  project   = "example"
  ref       = "main"
  path      = "ExampleSubFolder"
  recursive = true
}
