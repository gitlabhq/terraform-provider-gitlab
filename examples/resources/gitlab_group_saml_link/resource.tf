resource "gitlab_group_saml_link" "test" {
  group_id        = "12345"
  access_level    = "developer"
  saml_group_name = "samlgroupname1"
}
