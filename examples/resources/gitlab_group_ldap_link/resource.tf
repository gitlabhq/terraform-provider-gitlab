resource "gitlab_group_ldap_link" "test" {
  group_id      = "12345"
  cn            = "testuser"
  group_access  = "developer"
  ldap_provider = "ldapmain"
}
