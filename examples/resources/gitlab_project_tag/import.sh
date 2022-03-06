# Gitlab project tags can be imported with a key composed of `<project_id>:<tag_name>`, e.g.
terraform import gitlab_project_tag.example "12345:develop"

# NOTE: the `ref` attribute won't be available for imported `gitlab_project_tag` resources.
