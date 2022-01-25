resource "gitlab_project" "foo" {
  name = "foo-project"
}

resource gitlab_project_cluster "bar" {
  project                       = "${gitlab_project.foo.id}"
  name                          = "bar-cluster"
  domain                        = "example.com"
  enabled                       = true
  kubernetes_api_url            = "https://124.124.124"
  kubernetes_token              = "some-token"
  kubernetes_ca_cert            = "some-cert"
  kubernetes_namespace          = "namespace"
  kubernetes_authorization_type = "rbac"
  environment_scope             = "*"
  management_project_id         = "123456"
}
