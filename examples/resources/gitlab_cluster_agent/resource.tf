resource "gitlab_cluster_agent" "example" {
  project = "12345"
  name    = "agent-1"
}

// Optionally, configure the agent as described in
// https://docs.gitlab.com/ee/user/clusters/agent/install/index.html#create-an-agent-configuration-file
resource "gitlab_repository_file" "example_agent_config" {
  project        = gitlab_cluster_agent.example.project
  branch         = "main" // or use the `default_branch` attribute from a project data source / resource
  file_path      = ".gitlab/agents/${gitlab_cluster_agent.example.name}"
  content        = <<CONTENT
  gitops:
    ...
  CONTENT
  author_email   = "terraform@example.com"
  author_name    = "Terraform"
  commit_message = "feature: add agent config for ${gitlab_cluster_agent.example.name}"
}
