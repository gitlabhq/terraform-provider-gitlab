// Create token for an agent
resource "gitlab_cluster_agent_token" "example" {
  project     = "12345"
  agent_id    = 42
  name        = "some-token"
  description = "some token"
}

// The following example creates a GitLab Agent for Kubernetes in a given project,
// creates a token and install the `gitlab-agent` Helm Chart.
// (see https://gitlab.com/gitlab-org/charts/gitlab-agent)
data "gitlab_project" "this" {
  path_with_namespace = "my-org/example"
}

resource "gitlab_cluster_agent" "this" {
  project = data.gitlab_project.this
  name    = "my-agent"
}

resource "gitlab_cluster_agent_token" "this" {
  project     = data.gitlab_project.this
  agent_id    = gitlab_cluster_agent.this.id
  name        = "my-agent-token"
  description = "Token for the my-agent used with `gitlab-agent` Helm Chart"
}

resource "helm_release" "gitlab_agent" {
  name             = "gitlab-agent"
  namespace        = "gitlab-agent"
  create_namespace = true
  repository       = "https://charts.gitlab.io"
  chart            = "gitlab-agent"
  version          = "1.2.0"

  set {
    name  = "config.token"
    value = gitlab_cluster_agent_token.this.token
  }
}
