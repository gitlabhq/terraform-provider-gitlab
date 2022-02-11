resource "gitlab_user_custom_attribute" "attr" {
  user  = "42"
  key   = "location"
  value = "Greenland"
}
