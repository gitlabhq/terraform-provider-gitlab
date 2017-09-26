# 0.2.0 (Unreleased)

BACKWARDS INCOMPATIBILITIES:

Previous release of the terraform provider for gitlab used the API v3. This
release uses the v4 API. It means that the provider no longer works for
gitlab < 9.0.
If you specified a custom gitlab URL, make sure you switch to the `/api/v4` url.

IMPROVEMENTS:

* Block projects deletion until they are really deleted from server (#19)

NEW FEATURES:

* Add `cacert_file` and `insecure` options to the provider (#5)
* Move to the v4 api (#20)
* Support custom path for `gitlab_projects` (#21)
* New resource: `gitlab_labels` (#22)
* New resource: `gitlab_user` (#23)

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
