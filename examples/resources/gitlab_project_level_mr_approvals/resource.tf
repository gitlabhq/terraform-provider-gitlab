resource "gitlab_project" "foo" {
  name              = "Example"
  description       = "My example project"
}

resource "gitlab_project_level_mr_approvals" "foo" {
  project_id                                     = gitlab_project.foo.id
  reset_approvals_on_push                        = true
  disable_overriding_approvers_per_merge_request = false
  merge_requests_author_approval                 = false
  merge_requests_disable_committers_approval     = true
}
