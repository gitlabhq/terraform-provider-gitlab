resource "gitlab_project_custom_attribute" "attr" {
  project = "42"
  key     = "location"
  value   = "Greenland"
}
