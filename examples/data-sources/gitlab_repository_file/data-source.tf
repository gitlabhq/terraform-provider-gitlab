data "gitlab_repository_file" "example" {
  project   = "example"
  ref       = "main"
  file_path = "README.md"
}
