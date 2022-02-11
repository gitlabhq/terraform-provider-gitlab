resource "gitlab_group" "this" {
  name        = "example"
  path        = "example"
  description = "An example group"
}
resource "gitlab_project" "this" {
  name                   = "example"
  namespace_id           = gitlab_group.this.id
  initialize_with_readme = true
}
resource "gitlab_repository_file" "this" {
  project        = gitlab_project.this.id
  file_path      = "meow.txt"
  branch         = "main"
  content        = base64encode("Meow goes the cat")
  author_email   = "terraform@example.com"
  author_name    = "Terraform"
  commit_message = "feature: add meow file"
}
