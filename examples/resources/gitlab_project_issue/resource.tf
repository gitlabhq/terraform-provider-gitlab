resource "gitlab_project" "foo" {
  name             = "example project"
  description      = "Lorem Ipsum"
  visibility_level = "public"
}

resource "gitlab_project_issue" "welcome_issue" {
  project           = gitlab_project.foo.id
  title             = "Welcome!"
  description       = <<EOT
  Welcome to the ${gitlab_project.foo.name} project!

  EOT
  discussion_locked = true
}

output "welcome_issue_web_url" {
  value = data.gitlab_project_issue.web_url
}
