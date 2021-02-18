## 3.5.0 (Feb 18, 2021)

FEATURES:

* Add resource for project freeze periods ([#516](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#516 ))  

ENHANCEMENTS:

* Update go version and go-gitlab version ([#523](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#523)) 
* Support additional attributes in `gitlab_project_hook` ([#525](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#525)) 
* Link badges in README to proper workflows ([#527](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#527)) 
* gitlab_project: Check each push rule individually ([#531](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#531))
* Allow `full_path` in addition to `id` in gitlab_project data source ([#532](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#532)) 
* Update test fixtures for better usability ([#535](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#535)) 
* Check for state change on user delete ([#539](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#539)) 
* Increase gitlab_project import timeout ([#536](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#536))
* Add optional mirror options ([#554](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#554)) 
* Remove vendor folder ([#546](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#546)) 
* Add dependabot config ([#558](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#558)) 
* Fix EE tests actually running against CE ([#564](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#564)) 
* Fix EE test mounting license as a directory ([#568](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#568)) 

BUG FIXES:

* fix deploy_token expiration ([#510](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#510)) 
* Fix group_share_group nil pointer reference ([#555](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#555)) 

## 3.4.0 (Jan 14, 2021)

FEATURES:

* Support sharing a group with another group ([#511](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#511)) 
* Support Project Mirroring ([#512](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#512))

## 3.3.0 (Nov 30, 2020)

FEATURES:

* Support instance level CI variables ([#389](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#389))

ENHANCEMENTS

*  Add the pages_access_level parameter ([#472](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#472))
*  Do not fail when project member does not exist ([#473](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#473))
* Make the runners_token on the project secret ([#474](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#474))
*  Fix nil pointer dereference importing gitlab_user ([#490](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#490))
* Fix unit and acceptance tests not running ([#495](https://github.com/gitlabhq/terraform-provider-gitlab/pull/#495))

## 3.2.0 (Nov 20, 2020)

FEATURES:

* Project Approval Rules ([#250](https://github.com/gitlabhq/terraform-provider-gitlab/pull/https://github.com/gitlabhq/terraform-provider-gitlab/pull/250))

ENHANCEMENTS

* Documentation for expires_at ([#482](https://github.com/gitlabhq/terraform-provider-gitlab/pull/482))
* Update set-env github action command ([484](https://github.com/gitlabhq/terraform-provider-gitlab/pull/484))

## 3.1.0 (Oct 16, 2020)

ENHANCEMENTS:

* Enable custom UserAgent ([#451](https://github.com/gitlabhq/terraform-provider-gitlab/pull/451))
* gitlab_project_mirror: Mark URL as sensitive ([#458](https://github.com/gitlabhq/terraform-provider-gitlab/pull/458))
* Remove old-style variable interpolation ([#456](https://github.com/gitlabhq/terraform-provider-gitlab/pull/456))

BUG FIXES:

* add pagination for ListPipelineSchedules ([#454](https://github.com/gitlabhq/terraform-provider-gitlab/pull/454))

## 3.0.0 (Sept 23, 2020)

BREAKING CHANGES:

* Resource `gitlab_project_push_rules` has been removed. You now instead specify project push rules using the `push_rules` attribute on the `gitlab_project` resource.
* The `shared_with_groups` attribute has been removed from the `gitlab_project` resource (but not the data source). You may use the `gitlab_project_share_group` resource instead.

NOTES:

* If you are using the `environment_scope` attribute of `gitlab_project_variable` to manage multiple variables with the same key, it is recommended to use GitLab 13.4+. See [this related GitLab issue](https://gitlab.com/gitlab-org/gitlab/-/issues/9912) for older versions.
* The ID format of the `gitlab_project_variable` resource changed. The upgrade should be automatic.
* The default value of the `gitlab_project_variable` resource's `environment_scope` attribute has changed from `0` to `*`.

FEATURES:

* **New Data Source:** `gitlab_group_membership` ([#264](https://github.com/gitlabhq/terraform-provider-gitlab/issues/264))
* **New Resource:** `gitlab_instance_cluster` ([#367](https://github.com/gitlabhq/terraform-provider-gitlab/issues/367))
* **New Resource:** `gitlab_project_level_mr_approvals` ([#356](https://github.com/gitlabhq/terraform-provider-gitlab/issues/356))
* **New Resource:** `gitlab_project_mirror` ([#358](https://github.com/gitlabhq/terraform-provider-gitlab/issues/358))
* **New Resource:** `gitlab_service_pipelines_email` ([#375](https://github.com/gitlabhq/terraform-provider-gitlab/issues/375))

ENHANCEMENTS:

* data-source/gitlab_project: New attributes `packages_enabled`, `path_with_namespace` and `push_rules` ([#405](https://github.com/gitlabhq/terraform-provider-gitlab/issues/405), [#403](https://github.com/gitlabhq/terraform-provider-gitlab/issues/403), [#422](https://github.com/gitlabhq/terraform-provider-gitlab/issues/422))
* resource/gitlab_branch_protection: New `code_owner_approval_required` attribute ([#380](https://github.com/gitlabhq/terraform-provider-gitlab/issues/380))
* resource/gitlab_project: New attributes `packages_enabled`, `path_with_namespace`, and `push_rules` ([#405](https://github.com/gitlabhq/terraform-provider-gitlab/issues/405), [#403](https://github.com/gitlabhq/terraform-provider-gitlab/issues/403), [#422](https://github.com/gitlabhq/terraform-provider-gitlab/issues/422))
* resource/gitlab_group: New attributes `share_with_group_lock`, `project_creation_level`, `auto_devops_enabled`, `emails_disabled`, `mentions_disabled`, `subgroup_creation_level`, `require_two_factor_authentication`, and `two_factor_grace_period` ([#362](https://github.com/gitlabhq/terraform-provider-gitlab/issues/362))
* resource/gitlab_group: Automatically detect removal ([#267](https://github.com/gitlabhq/terraform-provider-gitlab/issues/267))
* resource/gitlab_group_label: Can now be imported ([#339](https://github.com/gitlabhq/terraform-provider-gitlab/issues/339))
* resource/gitlab_project: New `import_url` attribute ([#381](https://github.com/gitlabhq/terraform-provider-gitlab/issues/381))
* resource/gitlab_project_push_rules: Can now be imported ([#360](https://github.com/gitlabhq/terraform-provider-gitlab/issues/360))
* resource/gitlab_project_variable: Better error message when a masked variable fails validation ([#371](https://github.com/gitlabhq/terraform-provider-gitlab/issues/371))
* resource/gitlab_project_variable: Automatically detect removal ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
* resource/gitlab_service_jira: Automatically detect removal ([#337](https://github.com/gitlabhq/terraform-provider-gitlab/issues/337))
* resource/gitlab_user: The `email` attribute can be changed without forcing recreation ([#261](https://github.com/gitlabhq/terraform-provider-gitlab/issues/261))
* resource/gitlab_user: Require either the `password` or `reset_password` attribute to be set ([#262](https://github.com/gitlabhq/terraform-provider-gitlab/issues/262))

BUG FIXES:

* resource/gitlab_pipeline_schedule: Fix a rare error during deletion ([#364](https://github.com/gitlabhq/terraform-provider-gitlab/issues/364))
* resource/gitlab_pipeline_schedule_variable: Fix a rare error during deletion ([#364](https://github.com/gitlabhq/terraform-provider-gitlab/issues/364))
* resource/gitlab_project: Fix the `default_branch` attribute changing to `null` after first apply ([#343](https://github.com/gitlabhq/terraform-provider-gitlab/issues/343))
* resource/gitlab_project_share_group: Fix the `access_level` attribute not updating ([#421](https://github.com/gitlabhq/terraform-provider-gitlab/issues/421))
* resource/gitlab_project_share_group: Fix the share not working if the project is also managed ([#421](https://github.com/gitlabhq/terraform-provider-gitlab/issues/421))
* resource/gitlab_project_variable: Fix inconsistent reads for variables with non-unique keys ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
* resource/gitlab_project_variable: Change the default `environment_scope` from `0` to `*` ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
* resource/gitlab_service_jira: Fix a rare state inconsistency problem during creation ([#363](https://github.com/gitlabhq/terraform-provider-gitlab/issues/363))
* resource/gitlab_user: Fix some attributes saving incorrectly in state ([#261](https://github.com/gitlabhq/terraform-provider-gitlab/issues/261))

## 2.11.0 (July 24, 2020)

ENHANCEMENTS:
* Improvements to resource `gitlab_user` import
  ([#340](https://github.com/gitlabhq/terraform-provider-gitlab/issues/340))

## 2.10.0 (June 09, 2020)

FEATURES:
* **New Resource:** `gitlab_service_github`
  ([#311](https://github.com/gitlabhq/terraform-provider-gitlab/issues/311))

ENHANCEMENTS:
* add attribute `remove_source_branch_after_merge` to projects
  ([#289](https://github.com/gitlabhq/terraform-provider-gitlab/issues/289))

BUGFIXES:
* fix for flaky `gitlab_group` tests
  ([#320](https://github.com/gitlabhq/terraform-provider-gitlab/issues/320))
* Creating custom skip function for group_ldap_link tests.
  ([#328](https://github.com/gitlabhq/terraform-provider-gitlab/issues/328))

## 2.9.0 (June 01, 2020)

FEATURES:
* **New DataSource:** `gitlab_projects`
  ([#279](https://github.com/gitlabhq/terraform-provider-gitlab/issues/279))
* **New Resource:** `gitlab_deploy_token`
  ([#284](https://github.com/gitlabhq/terraform-provider-gitlab/issues/284))

ENHANCEMENTS:
* Add `management_project_id` for Group and Project Clusters
  ([#301](https://github.com/gitlabhq/terraform-provider-gitlab/issues/301))

## 2.8.0 (May 28, 2020)

FEATURES:
* **New Resource:** `gitlab_group_ldap_link`
  ([#296](https://github.com/gitlabhq/terraform-provider-gitlab/issues/296),
   [#316](https://github.com/gitlabhq/terraform-provider-gitlab/issues/316))

ENHANCEMENTS:

* Update resource gitlab_group_label to read labels from all pages
  ([#302](https://github.com/gitlabhq/terraform-provider-gitlab/issues/302))
* Provide a way to specify client cert and key
  ([#315](https://github.com/gitlabhq/terraform-provider-gitlab/issues/315))

BUGFIXES:
* Increase MaxIdleConnsPerHost in http.Transport
  ([#305](https://github.com/gitlabhq/terraform-provider-gitlab/issues/305))

## 2.7.0 (May 20, 2020)

* Implement `masked` parameters for `gitlab_group_variable`
  ([#271](https://github.com/gitlabhq/terraform-provider-gitlab/issues/271))

## 2.6.0 (April 08, 2020)

ENHANCEMENTS:
* Add jira flags
  ([#274](https://github.com/gitlabhq/terraform-provider-gitlab/issues/274))

## 2.5.1 (April 06, 2020)

BUGFIXES:

* Support for soft-delete of groups and projects in Gitlab Enterprise Edition 
  ([#282](https://github.com/gitlabhq/terraform-provider-gitlab/issues/282),
   [#283](https://github.com/gitlabhq/terraform-provider-gitlab/issues/283),
   [#285](https://github.com/gitlabhq/terraform-provider-gitlab/issues/285),
   [#291](https://github.com/gitlabhq/terraform-provider-gitlab/issues/291))

ENHANCEMENTS:
* Switched from Travis CI to Github Actions 
  ([#216](https://github.com/gitlabhq/terraform-provider-gitlab/issues/216))

## 2.5.0 (December 05, 2019)

ENHANCEMENTS:
* Implement `lfs_enabled`, `request_access_enabled`, and `pipelines_enabled` parameters for `gitlab_project`
  ([#225](https://github.com/gitlabhq/terraform-provider-gitlab/pull/225),
   [#226](https://github.com/gitlabhq/terraform-provider-gitlab/pull/226),
   [#227](https://github.com/gitlabhq/terraform-provider-gitlab/pull/227))

BUGFIXES:
* Fix label support when there is more than 20 labels on a project
  ([#229](https://github.com/gitlabhq/terraform-provider-gitlab/pull/229))
* Enable `environment_scope` for `gitlab_project_variable` lookup
  ([#228](https://github.com/gitlabhq/terraform-provider-gitlab/pull/229))
* Fix users data source when there is more than 20 users returned
  ([#230](https://github.com/gitlabhq/terraform-provider-gitlab/pull/230))

## 2.4.0 (November 28, 2019)

FEATURES:
* **New Resource:** `gitlab_group_label` ([#186](https://github.com/gitlabhq/terraform-provider-gitlab/pull/186))
* **New Resource:** `gitlab_group_cluster`
  ([#178](https://github.com/gitlabhq/terraform-provider-gitlab/pull/178))
* **New Resource:** `gitlab_pipeline_schedule_variable`
  ([#204](https://github.com/gitlabhq/terraform-provider-gitlab/pull/204))

ENHANCEMENTS:
* Add `runners_token` to gitlab groups ([#218](https://github.com/gitlabhq/terraform-provider-gitlab/pull/218))
* Add `reset_password` to `gitlab_user` ([#127](https://github.com/gitlabhq/terraform-provider-gitlab/pull/127))
* Update `access_level` available values ([#220](https://github.com/gitlabhq/terraform-provider-gitlab/pull/220))
* Make read callbacks graceful for `gitlab_project_share_group`, `gitlab_branch_protection` and
  `gitlab_label` resources ([#223](https://github.com/gitlabhq/terraform-provider-gitlab/pull/223))


BUGFIXES:
* Fix state not being updated for `gitlab_branch_protection`
  ([#166](https://github.com/gitlabhq/terraform-provider-gitlab/pull/166))
* Set ForceNew for `gitlab_pipeline_schedule` `project`
  ([#203](https://github.com/gitlabhq/terraform-provider-gitlab/pull/203))

## 2.3.0 (October 17, 2019)

*We would like to thank Gitlab, which has provided us a EE license. This project
is now tested against Gitlab CE and Gitlab EE.*

FEATURES:
* **New Resource:** `gitlab_project_push_rules` ([#163](https://github.com/gitlabhq/terraform-provider-gitlab/pull/163))
* **New Resource:** `gitlab_deploy_key_enable` ([#176](https://github.com/gitlabhq/terraform-provider-gitlab/pull/176))
* **New Resource:** `gitlab_project_share_group` ([#167](https://github.com/gitlabhq/terraform-provider-gitlab/pull/167))

ENHANCEMENTS:
* Add `initialize_with_readme` to `gitlab_project` ([#179](https://github.com/gitlabhq/terraform-provider-gitlab/issues/179))
* Add support for more variable options ([#169](https://github.com/gitlabhq/terraform-provider-gitlab/issues/169))
* Documentation improvements ([#168](https://github.com/gitlabhq/terraform-provider-gitlab/issues/168), [#187](https://github.com/gitlabhq/terraform-provider-gitlab/issues/187), [#171](https://github.com/gitlabhq/terraform-provider-gitlab/issues/171))

BUGFIXES:
* Fix tag protection URL
  ([#156](https://github.com/gitlabhq/terraform-provider-gitlab/issues/156))
* Properly manage the default branch in a git repo
  ([#158](https://github.com/gitlabhq/terraform-provider-gitlab/issues/158))
* Resolve triggers pagination issue by calling `GetPipelineTrigger`
  ([#173](https://github.com/gitlabhq/terraform-provider-gitlab/issues/173))

## 2.2.0 (June 12, 2019)

FEATURES:
* **New Resource:** `gitlab_service_jira` ([#101](https://github.com/gitlabhq/terraform-provider-gitlab/pull/101))
* **New Resource:** `gitlab_pipeline_schedule` ([#116](https://github.com/gitlabhq/terraform-provider-gitlab/pull/116))

ENHANCEMENTS:
* Add `archived` argument to `gitlab_project` ([#148](https://github.com/gitlabhq/terraform-provider-gitlab/issues/148))
* Add `managed` argument to `gitlab_project_cluster` ([#137](https://github.com/gitlabhq/terraform-provider-gitlab/issues/137))

## 2.1.0 (May 29, 2019)

FEATURES:
* **New Datasource**: `gitlab_group` ([#129](https://github.com/gitlabhq/terraform-provider-gitlab/issues/129))


## 2.0.0 (May 23, 2019)

This is the first release to support Terraform 0.12.

BACKWARDS INCOMPATIBILITIES:
* **all**: Previous versions of this provider silently removed state from state when
  Gitlab returned an error 404. Now we error on this and you must reconciliate
  the state (e.g. `terraform state rm`). We have done this because we can not
  make the difference between permission denied and resources removed outside of
  terraform (gitlab returns 404 in both cases)
  ([#130](https://github.com/gitlabhq/terraform-provider-gitlab/pull/130))


FEATURES:
* **New Resource:** `gitlab_tag_protection` ([#125](https://github.com/gitlabhq/terraform-provider-gitlab/pull/125))


ENHANCEMENTS:
* Add `container_registry_enabled` argument to `gitlab_project` ([#115](https://github.com/gitlabhq/terraform-provider-gitlab/issues/115))
* Add `shared_runners_enabled` argument to `gitlab_project` ([#134](https://github.com/gitlabhq/terraform-provider-gitlab/issues/134) [#104](https://github.com/gitlabhq/terraform-provider-gitlab/issues/104))

## 1.3.0 (May 03, 2019)

FEATURES:
* **New Resource:** `gitlab_service_slack` ([#96](https://github.com/gitlabhq/terraform-provider-gitlab/issues/96))
* **New Resource:** `gitlab_branch_protection` ([#68](https://github.com/gitlabhq/terraform-provider-gitlab/issues/68))

ENHANCEMENTS:
* Support for request/response logging when >`DEBUG` severity is set ([#93](https://github.com/gitlabhq/terraform-provider-gitlab/issues/93))
* Datasource `gitlab_user` supports user_id, email lookup and return lots of new attributes ([#102](https://github.com/gitlabhq/terraform-provider-gitlab/issues/102))
* Resource `gitlab_deploy_key` can now be imported ([#197](https://github.com/gitlabhq/terraform-provider-gitlab/issues/97))
* Add `tags` attribute for `gitlab_project` ([#106](https://github.com/gitlabhq/terraform-provider-gitlab/pull/106))


BUGFIXES:
* Documentation fixes ([#108](https://github.com/gitlabhq/terraform-provider-gitlab/issues/108), [#113](https://github.com/gitlabhq/terraform-provider-gitlab/issues/113))

## 1.2.0 (February 19, 2019)

FEATURES:

* **New Datasource:** `gitlab_users` ([#79](https://github.com/gitlabhq/terraform-provider-gitlab/issues/79))
* **New Resource:** `gitlab_pipeline_trigger` ([#82](https://github.com/gitlabhq/terraform-provider-gitlab/issues/82))
* **New Resource:** `gitlab_project_cluster` ([#87](https://github.com/gitlabhq/terraform-provider-gitlab/issues/87))

ENHANCEMENTS:

* Supports "No one" and "maintainer" permissions ([#83](https://github.com/gitlabhq/terraform-provider-gitlab/issues/83))
* `gitlab_project.shared_with_groups` is now order-independent ([#86](https://github.com/gitlabhq/terraform-provider-gitlab/issues/86))
* add `merge_method`, `only_allow_merge_if_*`, `approvals_before_merge` parameters to `gitlab_project` ([#72](https://github.com/gitlabhq/terraform-provider-gitlab/issues/72), [#88](https://github.com/gitlabhq/terraform-provider-gitlab/issues/88))


## 1.1.0 (January 14, 2019)

FEATURES:

* **New Resource:** `gitlab_project_membership`
* **New Resource:** `gitlab_group_membership` ([#8](https://github.com/gitlabhq/terraform-provider-gitlab/issues/8))
* **New Resource:** `gitlab_project_variable` ([#47](https://github.com/gitlabhq/terraform-provider-gitlab/issues/47))
* **New Resource:** `gitlab_group_variable` ([#47](https://github.com/gitlabhq/terraform-provider-gitlab/issues/47))

BACKWARDS INCOMPATIBILITIES:

`gitlab_project_membership` is not compatible with a previous *unreleased* version due to an id change resource will need to be reimported manually
e.g
```bash
terraform state rm gitlab_project_membership.foo
terraform import gitlab_project_membership.foo 12345:1337
```

## 1.0.0 (October 06, 2017)

BACKWARDS INCOMPATIBILITIES:

* This provider now uses the v4 api. It means that if you set up a custom API url, you need to update it to use the /api/v4 url. As a side effect, we no longer support Gitlab < 9.0. ([#20](https://github.com/gitlabhq/terraform-provider-gitlab/issues/20))
* We now support Parent ID for `gitlab_groups`. However, due to a limitation in
  the gitlab API, changing a Parent ID requires destroying and recreating the
  group. Since previous versions of this provider did not support it, there are
  chances that terraform will try do delete all your nested group when you
  update to 1.0.0. A workaround to prevent this is to use the `ignore_changes`
  lifecycle parameter. ([#28](https://github.com/gitlabhq/terraform-provider-gitlab/issues/28))

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

* **New Resource:** `gitlab_user` ([#23](https://github.com/gitlabhq/terraform-provider-gitlab/issues/23))
* **New Resource:** `gitlab_label` ([#22](https://github.com/gitlabhq/terraform-provider-gitlab/issues/22))

IMPROVEMENTS:

* Add `cacert_file` and `insecure` options to the provider. ([#5](https://github.com/gitlabhq/terraform-provider-gitlab/issues/5))
* Fix race conditions with `gitlab_project` deletion. ([#19](https://github.com/gitlabhq/terraform-provider-gitlab/issues/19))
* Add `parent_id` argument to `gitlab_group`. ([#28](https://github.com/gitlabhq/terraform-provider-gitlab/issues/28))
* Add support for `gitlab_project` import. ([#30](https://github.com/gitlabhq/terraform-provider-gitlab/issues/30))
* Add support for `gitlab_groups` import. ([#31](https://github.com/gitlabhq/terraform-provider-gitlab/issues/31))
* Add `path` argument for `gitlab_project`. ([#21](https://github.com/gitlabhq/terraform-provider-gitlab/issues/21))
* Fix indempotency issue with `gitlab_deploy_key` and white spaces. ([#34](https://github.com/gitlabhq/terraform-provider-gitlab/issues/34))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
