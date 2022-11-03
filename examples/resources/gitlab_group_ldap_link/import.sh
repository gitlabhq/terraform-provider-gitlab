# GitLab group ldap links can be imported using an id made up of
# `group_id:ldap_provider:cn`
terraform import gitlab_group_ldap_link.cn "12345:ldapmain:testuser"
# or `group_id:ldap_provider:filter
terraform import gitlab_group_ldap_link.test "12345:ldapmain:(objectClass=*)"
