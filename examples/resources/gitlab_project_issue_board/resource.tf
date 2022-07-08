resource "gitlab_project" "example" {
  name             = "example project"
  description      = "Lorem Ipsum"
  visibility_level = "public"
}

resource "gitlab_user" "example" {
  name     = "example"
  username = "example"
  email    = "example@example.com"
  password = "example1$$$"
}

resource "gitlab_project_membership" "example" {
  project_id   = gitlab_project.example.id
  user_id      = gitlab_user.example.id
  access_level = "developer"
}

resource "gitlab_project_milestone" "example" {
  project = gitlab_project.example.id
  title   = "m1"
}

resource "gitlab_project_issue_board" "this" {
  project = gitlab_project.example.id
  name    = "Test Issue Board"

  lists {
    assignee_id = gitlab_user.example.id
  }

  lists {
    milestone_id = gitlab_project_milestone.example.milestone_id
  }

  depends_on = [
    gitlab_project_membership.example
  ]
}

resource "gitlab_project_issue_board" "list_syntax" {
  project = gitlab_project.example.id
  name    = "Test Issue Board with list syntax"

  lists = [
    {
      assignee_id = gitlab_user.example.id
    },
    {
      milestone_id = gitlab_project_milestone.example.milestone_id
    }
  ]

  depends_on = [
    gitlab_project_membership.example
  ]
}
