---
layout: "gitlab"
page_title: "GitLab: gitlab_project_cluster"
sidebar_current: "docs-gitlab-resource-project_cluster"
description: |-
  Creates and manages project clusters for GitLab projects
---

# gitlab\_project\_cluster

This resource allows you to create and manage project clusters for your GitLab projects.
For further information on clusters, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/project/clusters/index.html).


## Example Usage

```hcl
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
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The id of the project to add the cluster to.

* `name` - (Required, string) The name of cluster.

* `domain` - (Optional, string) The base domain of the cluster.

* `enabled` - (Optional, boolean) Determines if cluster is active or not. Defaults to `true`. This attribute cannot be read.

* `managed` - (Optional, boolean) Determines if cluster is managed by gitlab or not. Defaults to `true`. This attribute cannot be read.

* `kubernetes_api_url` - (Required, string) The URL to access the Kubernetes API.

* `kubernetes_token` - (Required, string) The token to authenticate against Kubernetes.

* `kubernetes_ca_cert` - (Optional, string) TLS certificate (needed if API is using a self-signed TLS certificate).

* `kubernetes_namespace` - (Optional, string) The unique namespace related to the project.

* `kubernetes_authorization_type` - (Optional, string) The cluster authorization type. Valid values are `rbac`, `abac`, `unknown_authorization`. Defaults to `rbac`.

* `environment_scope` - (Optional, string) The associated environment to the cluster. Defaults to `*`.

## Import

GitLab project clusters can be imported using an id made up of `projectid:clusterid`, e.g.

```
$ terraform import gitlab_project_cluster.bar 123:321
```
