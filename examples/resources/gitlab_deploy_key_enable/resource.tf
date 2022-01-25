# A repo to host the deployment key
resource "gitlab_project" "parent" {
  name = "parent_project"
}

# A second repo to use the deployment key from the parent project
resource "gitlab_project" "foo" {
  name = "foo_project"
}

# Upload a deployment key for the parent repo
resource "gitlab_deploy_key" "parent" {
  project = "${gitlab_project.parent.id}"
  title = "Example deploy key"
  key = "ssh-rsa AAAA..."
}

# Enable the deployment key on the second repo
resource "gitlab_deploy_key_enable" "foo" {
  project = "${gitlab_project.foo.id}"
  key_id = "${gitlab_deploy_key.parent.id}"
}
