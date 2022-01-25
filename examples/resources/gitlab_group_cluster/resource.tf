resource "gitlab_group" "foo" {
  name = "foo-group"
  path = "foo-path"
}

resource "gitlab_group_cluster" "bar" {
  group                       = "${gitlab_group.foo.id}"
  name                          = "bar-cluster"
  domain                        = "example.com"
  enabled                       = true
  kubernetes_api_url            = "https://124.124.124"
  kubernetes_token              = "some-token"
  kubernetes_ca_cert            = "some-cert"
  kubernetes_authorization_type = "rbac"
  environment_scope             = "*"
  management_project_id         = "123456"
}
