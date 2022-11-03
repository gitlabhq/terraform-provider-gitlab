resource "gitlab_group_ldap_link" "cn" {
  group_id      = "12345"
  cn            = "testuser"
  group_access  = "developer"
  ldap_provider = "ldapmain"
}

resource "gitlab_group_ldap_link" "filter" {
  group_id      = "12345"
  filter        = "(objectClass=*)"
  group_access  = "reporter"
  ldap_provider = "ldapmain"
}