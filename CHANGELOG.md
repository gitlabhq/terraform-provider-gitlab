## 2.2.1 (Unreleased)
## 2.2.0 (June 12, 2019)

FEATURES:
* **New Resource:** `gitlab_service_jira` ([#101](https://github.com/terraform-providers/terraform-provider-gitlab/pull/101))
* **New Resource:** `gitlab_pipeline_schedule` ([#116](https://github.com/terraform-providers/terraform-provider-gitlab/pull/116))

ENHANCEMENTS:
* Add `archived` argument to `gitlab_project` ([#148](https://github.com/terraform-providers/terraform-provider-gitlab/issues/148))
* Add `managed` argument to `gitlab_project_cluster` ([#137](https://github.com/terraform-providers/terraform-provider-gitlab/issues/137))

## 2.1.0 (May 29, 2019)

FEATURES:
* **New Datasource**: `gitlab_group` ([#129](https://github.com/terraform-providers/terraform-provider-gitlab/issues/129))


## 2.0.0 (May 23, 2019)

This is the first release to support Terraform 0.12.

BACKWARDS INCOMPATIBILITIES:
* **all**: Previous versions of this provider silently removed state from state when
  Gitlab returned an error 404. Now we error on this and you must reconciliate
  the state (e.g. `terraform state rm`). We have done this because we can not
  make the difference between permission denied and resources removed outside of
  terraform (gitlab returns 404 in both cases)
  ([#130](https://github.com/terraform-providers/terraform-provider-gitlab/pull/130))


FEATURES:
* **New Resource:** `gitlab_tag_protection` ([#125](https://github.com/terraform-providers/terraform-provider-gitlab/pull/125))


ENHANCEMENTS:
* Add `container_registry_enabled` argument to `gitlab_project` ([#115](https://github.com/terraform-providers/terraform-provider-gitlab/issues/115))
* Add `shared_runners_enabled` argument to `gitlab_project` ([#134](https://github.com/terraform-providers/terraform-provider-gitlab/issues/134) [#104](https://github.com/terraform-providers/terraform-provider-gitlab/issues/104))

## 1.3.0 (May 03, 2019)

FEATURES:
* **New Resource:** `gitlab_service_slack` ([#96](https://github.com/terraform-providers/terraform-provider-gitlab/issues/96))
* **New Resource:** `gitlab_branch_protection` ([#68](https://github.com/terraform-providers/terraform-provider-gitlab/issues/68))

ENHANCEMENTS:
* Support for request/response logging when >`DEBUG` severity is set ([#93](https://github.com/terraform-providers/terraform-provider-gitlab/issues/93))
* Datasource `gitlab_user` supports user_id, email lookup and return lots of new attributes ([#102](https://github.com/terraform-providers/terraform-provider-gitlab/issues/102))
* Resource `gitlab_deploy_key` can now be imported ([#197](https://github.com/terraform-providers/terraform-provider-gitlab/issues/97))
* Add `tags` attribute for `gitlab_project` ([#106](https://github.com/terraform-providers/terraform-provider-gitlab/pull/106))


BUGFIXES:
* Documentation fixes ([#108](https://github.com/terraform-providers/terraform-provider-gitlab/issues/108), [#113](https://github.com/terraform-providers/terraform-provider-gitlab/issues/113))

## 1.2.0 (February 19, 2019)

FEATURES:

* **New Datasource:** `gitlab_users` ([#79](https://github.com/terraform-providers/terraform-provider-gitlab/issues/79))
* **New Resource:** `gitlab_pipeline_trigger` ([#82](https://github.com/terraform-providers/terraform-provider-gitlab/issues/82))
* **New Resource:** `gitlab_project_cluster` ([#87](https://github.com/terraform-providers/terraform-provider-gitlab/issues/87))

ENHANCEMENTS:

* Supports "No one" and "maintainer" permissions ([#83](https://github.com/terraform-providers/terraform-provider-gitlab/issues/83))
* `gitlab_project.shared_with_groups` is now order-independent ([#86](https://github.com/terraform-providers/terraform-provider-gitlab/issues/86))
* add `merge_method`, `only_allow_merge_if_*`, `approvals_before_merge` parameters to `gitlab_project` ([#72](https://github.com/terraform-providers/terraform-provider-gitlab/issues/72), [#88](https://github.com/terraform-providers/terraform-provider-gitlab/issues/88))


## 1.1.0 (January 14, 2019)

FEATURES:

* **New Resource:** `gitlab_project_membership`
* **New Resource:** `gitlab_group_membership` ([#8](https://github.com/terraform-providers/terraform-provider-gitlab/issues/8))
* **New Resource:** `gitlab_project_variable` ([#47](https://github.com/terraform-providers/terraform-provider-gitlab/issues/47))
* **New Resource:** `gitlab_group_variable` ([#47](https://github.com/terraform-providers/terraform-provider-gitlab/issues/47))

BACKWARDS INCOMPATIBILITIES:

`gitlab_project_membership` is not compatible with a previous *unreleased* version due to an id change resource will need to be reimported manually
e.g
```bash
terraform state rm gitlab_project_membership.foo
terraform import gitlab_project_membership.foo 12345:1337
```

## 1.0.0 (October 06, 2017)

BACKWARDS INCOMPATIBILITIES:

* This provider now uses the v4 api. It means that if you set up a custom API url, you need to update it to use the /api/v4 url. As a side effect, we no longer support Gitlab < 9.0. ([#20](https://github.com/terraform-providers/terraform-provider-gitlab/issues/20))
* We now support Parent ID for `gitlab_groups`. However, due to a limitation in
  the gitlab API, changing a Parent ID requires destroying and recreating the
  group. Since previous versions of this provider did not support it, there are
  chances that terraform will try do delete all your nested group when you
  update to 1.0.0. A workaround to prevent this is to use the `ignore_changes`
  lifecycle parameter. ([#28](https://github.com/terraform-providers/terraform-provider-gitlab/issues/28))

```
resource "gitlab_group" "nested_group" {
  name = "bar-name-%d"
  path = "bar-path-%d"
  lifecycle {
    ignore_changes = ["parent_id"]
  }
}
```

FEATURES:

* **New Resource:** `gitlab_user` ([#23](https://github.com/terraform-providers/terraform-provider-gitlab/issues/23))
* **New Resource:** `gitlab_label` ([#22](https://github.com/terraform-providers/terraform-provider-gitlab/issues/22))

IMPROVEMENTS:

* Add `cacert_file` and `insecure` options to the provider. ([#5](https://github.com/terraform-providers/terraform-provider-gitlab/issues/5))
* Fix race conditions with `gitlab_project` deletion. ([#19](https://github.com/terraform-providers/terraform-provider-gitlab/issues/19))
* Add `parent_id` argument to `gitlab_group`. ([#28](https://github.com/terraform-providers/terraform-provider-gitlab/issues/28))
* Add support for `gitlab_project` import. ([#30](https://github.com/terraform-providers/terraform-provider-gitlab/issues/30))
* Add support for `gitlab_groups` import. ([#31](https://github.com/terraform-providers/terraform-provider-gitlab/issues/31))
* Add `path` argument for `gitlab_project`. ([#21](https://github.com/terraform-providers/terraform-provider-gitlab/issues/21))
* Fix indempotency issue with `gitlab_deploy_key` and white spaces. ([#34](https://github.com/terraform-providers/terraform-provider-gitlab/issues/34))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
