---
page_title: "Terraform GitLab Provider Version 16.0 Upgrade Guide"
subcategory: "Upgrade Guides"
---

# Upgrade to Terraform GitLab Provider Version 16.0

The GitLab 16.0 major milestone introduced a couple of breaking changes that this
release addresses. In addition, the provider itself breaks a bunch of interfaces 
which may require actions on the users side. These are described below:

## Change of resource id formats

**Affected resources**:

- `gitlab_deploy_token`
- `gitlab_deploy_key`
- `gitlab_project_hook`
- `gitlab_group_label`
- `gitlab_project_label`
- `gitlab_pipeline_schedule_variable`
- `gitlab_group_ldap_link`
- `gitlab_pipeline_trigger`
- `gitlab_pipeline_schedule`

Some resource `id` formats weren't uniquely identifiable and did not contain
all the information to retrieve a particular resource from the GitLab API
given this id.

Therefore, you may need to change your `terraform import` commands to use fully
identifiable resource ids as described in the respective docs.
This may also affect any use of the `id` attribute (`gitlab_*.*.id`) access in
your Terraform configuration code.

For `gitlab_deploy_key` and `gitlab_pipeline_schedule` you need
to use the `deploy_key_id` and `pipeline_schedule_id` attributes, respectively to
refer to them in sub-resources, like `gitlab_deploy_key_enable` and `gitlab_pipeline_schedule_variable`.

## Change variable `value` attribute to non-sensitive

**Affected resources**:

- `gitlab_instance_variable`
- `gitlab_group_variable`
- `gitlab_project_variable`

The `value` attribute of the `gitlab_*_variable` resources has been changed
from `sensitive` to non-sensitive.

Therefore, you may want to use the `sensitive()` or `nonsensitive()` Terraform functions.

## Require `expires_at` attribute for Project Access Tokens

The `expires_at` attribate for the `gitlab_project_access_token` resource is required in
version 16.0.0.

In 16.0.1 and later, this attribute is optional again due to GitLab applying a default 
when it's empty, however setting it to a date too far in the future may cause an error
depending on the configuration of your GitLab instance.

## Change `project_id` attribute to `project`

**Affected resources**:

- `gitlab_project_freeze_period`
- `gitlab_project_level_mr_approvals`
- `gitlab_project_membership`
- `gitlab_project_share_group`

Some resources used a `project_id` attribute to identify a project by numerical id.
This attribute has been removed in favor of a new `project` attribute which supports
both numerical ids and full paths to the project to identify id.
This aligns with the rest of the project-scoped resources.

## Change `group_id` attribute to `group`

- `gitlab_group_ldap_link`

Some resources used a `group_id` attribute to identify a group by numerical id.
This attribute has been removed in favor of a new `group` attribute which supports
both numerical ids and full paths to the group to identify id.
This aligns with the rest of the group-scoped resources.

## Deprecate `gitlab_service_*` resources

All the `gitlab_service_*` resources have been deprecated in favor
of the new `gitlab_integration_*` resources.
Make sure to adapt to the new ones within the next 3 releases as we'll be
removing the `gitlab_service_*` resources with the upcoming 16.3 release.

## Deprecate `gitlab_label` resource

The `gitlab_label` resource has been deprecated in favor of the new
`gitlab_project_label` resource.
Make sure to adapt to the new resource within the next 3 releases as we'll be
removing the `gitlab_label` resource with the upcoming 16.3 release.

## Remove support for unencoded test in `gitlab_repository_file`

Support for non-base64 encoded text in `gitlab_repository_file` has been removed.
If unencoded values are used, terraform will now return an error noting 
`Invalid base64 string in "content"`. 
Instead, use the [`base64encode()`](https://developer.hashicorp.com/terraform/language/functions/base64encode)
function from terraform to encode any values if they are not already encoded.

## Misc removals

- The `gitlab_managed_license` resource has been removed. There is no longer an upstream GitLab API for it.
- The `operations_access_level` attribute was removed from the `gitlab_project` resource and data sources.