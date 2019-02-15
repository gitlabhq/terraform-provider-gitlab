## 1.2.0 (Unreleased)

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
