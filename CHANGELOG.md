## 0.2.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES:

* This provider now uses the v4 api. It means that if you set up a custom API url, you need to update it to use the /api/v4 url. As a side effect, we no longer support Gitlab < 9.0. [GH-20]

IMPROVEMENTS:

* Add `cacert_file` and `insecure` options to the provider. [GH-5]
* Fix race conditions with `gitlab_projects` deletion. [GH-19]

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
