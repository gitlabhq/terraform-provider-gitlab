resource "gitlab_group_ldap_link" "test" {
  group_id = "12345"
  cn = "testuser"
  access_level = "developer"
  ldap_provider = "ldapmain"
}
