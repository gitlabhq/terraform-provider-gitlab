data "gitlab_users" "example" {
  sort = "desc"
  order_by = "name"
  created_before = "2019-01-01"
}

data "gitlab_users" "example-two" {
  search = "username"
}
