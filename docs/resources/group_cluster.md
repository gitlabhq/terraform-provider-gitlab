# gitlab\_group\_cluster

This resource allows you to create and manage group clusters for your GitLab groups.
For further information on clusters, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/group/clusters/index.html).

## Example Usage

```hcl
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
```

## Argument Reference

The following arguments are supported:

* `group` - (Required, string) The id of the group to add the cluster to.

* `name` - (Required, string) The name of cluster.

* `domain` - (Optional, string) The base domain of the cluster.

* `enabled` - (Optional, boolean) Determines if cluster is active or not. Defaults to `true`. This attribute cannot be read.

* `managed` - (Optional, boolean) Determines if cluster is managed by gitlab or not. Defaults to `true`. This attribute cannot be read.

* `kubernetes_api_url` - (Required, string) The URL to access the Kubernetes API.

* `kubernetes_token` - (Required, string) The token to authenticate against Kubernetes.

* `kubernetes_ca_cert` - (Optional, string) TLS certificate (needed if API is using a self-signed TLS certificate).

* `kubernetes_authorization_type` - (Optional, string) The cluster authorization type. Valid values are `rbac`, `abac`, `unknown_authorization`. Defaults to `rbac`.

* `environment_scope` - (Optional, string) The associated environment to the cluster. Defaults to `*`.

* `management_project_id` - (Optional, string) The ID of the management project for the cluster.

## Import

GitLab group clusters can be imported using an id made up of `groupid:clusterid`, e.g.

```
$ terraform import gitlab_group_cluster.bar 123:321
```
