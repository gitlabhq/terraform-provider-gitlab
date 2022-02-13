data "gitlab_instance_deploy_keys" "example" {}

# only public deploy keys
data "gitlab_instance_deploy_keys" "example" {
  public = true
}
