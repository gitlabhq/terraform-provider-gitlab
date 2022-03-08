data "gitlab_project" "foo" {
  id = "foo/bar/baz"
}

data "gitlab_project_issue" "welcome_issue" {
  project = data.gitlab_project.foo.id
  iid     = 1
}

output "welcome_issue_web_url" {
  value = data.gitlab_project_issue.web_url
}
