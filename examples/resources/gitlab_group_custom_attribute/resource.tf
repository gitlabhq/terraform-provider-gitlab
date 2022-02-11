resource "gitlab_group_custom_attribute" "attr" {
  group = "42"
  key   = "location"
  value = "Greenland"
}
