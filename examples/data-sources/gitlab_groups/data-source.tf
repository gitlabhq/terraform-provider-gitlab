data "gitlab_groups" "example" {
  sort     = "desc"
  order_by = "name"
}

data "gitlab_groups" "example-two" {
  search = "GitLab"
}
