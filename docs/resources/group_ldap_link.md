# gitlab\_group\_ldap\_link

This resource allows you to add an LDAP link to an existing GitLab group.

## Example Usage

```hcl
resource "gitlab_group_ldap_link" "test" {
  group_id = "12345"
  cn = "testuser"
  access_level = "developer"
  ldap_provider = "ldapmain"
}
```

## Argument Reference

The following arguments are supported:

* `group_id` - (Required) The id of the GitLab group.

* `cn` - (Required) The CN of the LDAP group to link with.

* `access_level` - (Required) Acceptable values are: guest, reporter, developer, maintainer, owner.

* `ldap_provider` - (Required) The name of the LDAP provider as stored in the GitLab database.

## Import

GitLab group ldap links can be imported using an id made up of `ldap_provider:cn`, e.g.

```
$ terraform import gitlab_group_ldap_link.test "ldapmain:testuser"
```
