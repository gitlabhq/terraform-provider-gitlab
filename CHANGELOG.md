## 1.0.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES:

* This provider now uses the v4 api. It means that if you set up a custom API url, you need to update it to use the /api/v4 url. As a side effect, we no longer support Gitlab < 9.0. [GH-20]
* We now support Parent ID for `gitlab_groups`. However, due to a limitation in
  the gitlab API, changing a Parent ID requires destroying and recreating the
  group. Since previous versions of this provider did not support it, there are
  chances that terraform will try do delete all your nested group when you
  update to 0.2.0. A workaround to prevent this is to use the `ignore_changes`
  lifecycle parameter. [GH-28]

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

* **New Resource:** `gitlab_user` [GH-23]
* **New Resource:** `gitlab_label` [GH-22]

IMPROVEMENTS:

* Add `cacert_file` and `insecure` options to the provider. [GH-5]
* Fix race conditions with `gitlab_project` deletion. [GH-19]
* Add `parent_id` argument to `gitlab_group`. [GH-28]
* Add support for `gitlab_project` import. [GH-30]
* Add support for `gitlab_groups` import. [GH-31]
* Add `path` argument for `gitlab_project`. [GH-21]
* Fix indempotency issue with `gitlab_deploy_key` and white spaces. [GH-34]

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
