## 17.7.0 (2024-12-19)

This release was tested against GitLab 17.7, 17.6, and 17.5 for both CE and EE

### FEATURES (1 change)

- resource/gitlab_instance_service_account:  [Allows managing instance-level service accounts](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/0da40eceb4272d3fb856139198da777d08960711) by @ipsavitsky234 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2220)) 

### IMPROVEMENTS (2 changes)

- resource/gitlab_integration_jira:  [Align Jira integration resource with the API](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/e3b04f9babac66bc96244d32d79c304e639161cd) by @slaup1 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2221)) 
- resource/gitlab_project:  [Add support for 'prevent_merge_without_jira_issue'](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/6da24d2c30b017d00f766a6bdb9deb31877be242) by @slaup1 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2223)) 
- resource/gitlab_branch: [Add support for `keep_on_destroy`, which will preserve a branch when a `destroy` operation is performed](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/3fa3aae0ea19a71354bc7af13196bac4bd7fdd9a) by @jae-sorcero ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2192))

### BUG FIXES (4 changes)

- resource/gitlab_group_service_account_access_token:  [Fixed an issue where revoking a service account token required more permissions than the documentation suggested when using SaaS](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b4814f9bec3cbe904cc747f98f1d4adcff3691d0) by @u120230 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2240)) 
- resource/gitlab_group:  [Fixed an issue where older versions of GitLab may fail due to missing `DefaultBranchProtectionDefaults`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/c12242b1c7209cd00db546817a94c72e62ff4d6f) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2228))
- resource/gitlab_application_settings:  [Fixed an issue where `DefaultBranchProtectionDefaults` didn't apply properly](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/0f8c9e14e92adbee43087c483ccc6ecdfafcc35b) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2229))
- resource/gitlab_project: [Fixed an issue where `restrict_user_defined_variables` caused a perpetual plan when not included in the config](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/9a371bb386650a1efd79c6b3d8d9073a9d55e81d) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2243))

## 17.6.1 (2024-11-27)

This release was tested against GitLab 17.6, 17.5, and 17.4 for both CE and EE

### IMPROVEMENTS (1 change)

- resource/gitlab_group_service_account_token: [Add support for `rotation_configuration` with a limited set of attributes](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/7bad2a4bf0fd6e777c2d3598b35d46bc287135fc) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2206))

### BUG FIXES (2 changes)

- resource/gitlab_project_compliance_frameworks: [Fixed an issue where project compliance frameworks didn't delete compliance frameworks properly. Fixed a bug where compliance frameworks weren't set properly on a project where they were removed outside of TF](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/6fe104b21bf5f07aff62b9751fa3d2d66c5ba175) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2214))
- resource/gitlab_group: [Fixed an issue where an `allowed_email_domain_list` added outside of TF would be unset when upgrading the provider to 17.6](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/4b76021b5a0633054eccb5e0b24435cf8dd64a8c) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2217))

## 17.6.0 (2024-11-21)

This release was tested against GitLab 17.6, 17.5, and 17.4 for both CE and EE

### FEATURES (4 changes)

- **New Resource** resource/gitlab_project_compliance_frameworks: [Add new resource `gitlab_project_compliance_frameworks`, allowing managing multiple compliance frameworks on a single project](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/5022e5ee705a1ff926d2391e2e61a6ef3968db5b) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2209))
- **New Datasource** datasource/gitlab_group_billable_member_memberships:  [Add new data source `gitlab_group_billable_member_memberships`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/8cb70059c55bd121df6f84dd9c306b8559303754) by @mdroessler ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2144)) 
- **New Datasource** datasource/gitlab_project_environments:  [Add new data source `gitlab_project_environments`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/cca616faa31acc2702b9cd873a0e59b79aa00c06) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2172)) 
- **New Datasource** data/gitlab_runners:  [Add new data source `gitlab_runners`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b02e9ee2a0518a6e1b4133824661b9a6ebd221f7) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2164)) 

### IMPROVEMENTS (17 changes)

- docs:  [add import block examples to resources](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d350937502c4ee11c8c99ac89e7983bae03cc439) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2169)) 
- data/gitlab_project: [Add `ci_pipeline_variables_minimum_role` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/59215d220e064fc9d42d4f579031730e8aaec2bf) by @hfyngvason ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2208))
- data/gitlab_projects: [Add `ci_pipeline_variables_minimum_role` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/59215d220e064fc9d42d4f579031730e8aaec2bf) by @hfyngvason ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2208))
- data/gitlab_metadata: [Add `external_kas_proxy_url` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/ad66d449ab7e6257dc6787d2232cdf80c4b49a9e) by @timofurrer ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2203))
- resource/gitlab_user_runner:  [Add `maintenance_note` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/269371d4c91fc5ec333cd7b596964b63dfa60c70) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2170)) 
- resource/gitlab_group: [Add `allowed_email_domains_list` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/34024e7d01d1c8a11f3799416eba64c96653ba75) by @jeremad ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2204))
- resource/gitlab_project: [Add `ci_pipeline_variables_minimum_role` setting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/59215d220e064fc9d42d4f579031730e8aaec2bf) by @hfyngvason ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2208))
- resource/gitlab_project_environment: [Add support for multiple attributes related to cluster agents](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/ff6c1780c4f168aeb721be188098cdeda0336728) by @timofurrer ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2168))
- resource/gitlab_project_hook:  [Add `custom_headers` support](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/9a396394049d469e06afc43841d85b4fc807fb07) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2166)) 
- resource/gitlab_group_hook:  [Add `custom_headers` support](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/9a396394049d469e06afc43841d85b4fc807fb07) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2166)) 
- resource/gitlab_branch_protection:  [Add `deploy_key` support to `allowed_to_push`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/3d988071a93df18d93d9ead00a47bc4bf905fe45) by @nicolas.boussuge.fintech ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2157)) 
- resource/gitlab_group_security_policy_attachment:  [Update documentation to make usage with compliance frameworks and policies more clear](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/72457495cad307146cb1200318c447a1fb6bbe2a) by @odupre ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2210)) 
- resource/personal_access_token:  [Remove the "RequiresReplace()" modifier on rotation_config for access tokens](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/399f3c2a067f1e683258820ac88df7abd187e310) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2174)) 
- resource/project_access_token:  [Remove the "RequiresReplace()" modifier on rotation_config for access tokens](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/399f3c2a067f1e683258820ac88df7abd187e310) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2174)) 
- resource/group_access_token:  [Remove the "RequiresReplace()" modifier on rotation_config for access tokens](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/399f3c2a067f1e683258820ac88df7abd187e310) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2174)) 
- resource/gitlab_service_account_access_token:  [Improve example documentation and include better information on required permissions for gitlab.com](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/38b3c7e8b9ad27dd8b9ba452ceb6fd179801e550) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2201))
- resource/gitlab_member_role: [Add documentation examples and import examples](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/93ceb96da2be4e75c15fb5e552106c1aad5c5407) by @john.shimmin ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2186))


### BUG FIXES (8 changes)

- resource/gitlab_group_label:  [Fix an issue where `apply` fails when changing the `name` attribute](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b6dd163bb62c52c8cad432b92e638220687887a4) by @krzywykrzych ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2200)) 
- resource/gitlab_application_settings:  [Fix panic when using `valid_runner_registrars`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/55dc97d6ab2308f557595a4d50f3015776006507) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2173)) 
- resource/gitlab_service_account_access_token:  [Fix an issue where non-admin tokens cannot use the service account access token resource](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/20e8244e61544eceacd00be025b2bdd1a6299419) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2189))
- resource/gitlab_service_account_access_token:  [Fix invalid HCL in the example](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/a752175ba841041fd5a8889c072bae569dadaf1b) by @telliffoshea ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2167))
- resource/gitlab_member_role: [Fix an issue where `enabled_permissions` could cause issues due to sorting](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/93ceb96da2be4e75c15fb5e552106c1aad5c5407) by @john.shimmin ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2186))
- resource/gitlab_group_access_token: [Fix an issue where removing `rotation_configuration` after an `apply` operation caused a provider panic](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/40dbca0ff63b6c4639a806a156f201cf3fea202f) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2175))
- resource/gitlab_project_access_token: [Fix an issue where removing `rotation_configuration` after an `apply` operation caused a provider panic](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/40dbca0ff63b6c4639a806a156f201cf3fea202f) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2175))
- resource/gitlab_personal_access_token: [Fix an issue where removing `rotation_configuration` after an `apply` operation caused a provider panic](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/40dbca0ff63b6c4639a806a156f201cf3fea202f) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2175))

## 17.5.0 (2024-10-17)

This release was tested against GitLab 17.5, 17.4, and 17.3 for both CE and EE

### FEATURES (4 changes)

- **New Resource** resource/gitlab_group_service_account: [Allows creating service accounts for a gitlab group](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/a1e1686636ee2d901344163bac22bc032ad52056) by @michal.szczepek ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2122))
- **New Resource** resource/gitlab_group_service_account_token: [Allows creating service account tokens for a service account](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/a1e1686636ee2d901344163bac22bc032ad52056) by @michal.szczepek ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2122))
- **New Datasource** data/gitlab_group_service_account: [Allows creating service accounts for a gitlab group](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/a1e1686636ee2d901344163bac22bc032ad52056) by @michal.szczepek ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2122))
- **New Datasource** data/gitlab_group_provisioned_users: [Allows listing users provisioned by an Identity Provider](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/7daffab21eedd594746671a28becd20adba15f8c) by @nvh04121 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2107))

### IMPROVEMENTS (9 changes)

- provider/gitlab:  [Add support for `headers` to the provider to allow setting custom headers](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/72f1a49dddcaa9ad18a81a604066b9959f565bed) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2136))
- resource/gitlab_application_settings: [Add support for 20 additional new fields](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/5a6321fb096c80be4319dc29cf2186f66163bd88) by @Jitsusama ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2110))
- resource/gitlab_application_settings: [Add support for 23 (for a total of 43 in 17.5) additional new fields](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/4c07217976e4cebcbda7ee02bbda71a6696e9dac) by @Jitsusama ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2109))
- resource/gitlab_project:  [Support `allow_pipeline_trigger_approve_deployment`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/2497fc45f72577292412a398302b508c357e9a5c) by @loozhengyuan ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2141))
- resource/gitlab_project:  [Support `model_experiments_access_level` and `model_registry_access_level`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/aec6e45be6f6ee687062afa96b8b8b9347540bad) by @jeremad ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2125))
- resource/gitlab_group:  [Add support for specifying the default branch](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b0983be7c8cec2651af5ce82cdbda1caf5e82c56) by @jae-sorcero ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2155))
- resource/gitlab_group:  [Add `no one` as a valid permissions when using `default_branch_protection_defaults`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/bb0242b353c34907d88e47a418a51404228af7d9) by @jae-sorcero ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2147))
- resource/gitlab_personal_access_token:  [Add support for `rotation_configuration`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/afa261c798635521fd660b04b1a6a7d9100113a7) by @ruben.aleman ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2143))
- data/gitlab_pipeline_schedules:  [Add `id` for each schedule found using the datasource](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/169c4daf5bf51e9f9176ea76c00e518117e0de4e) by @nicolas.boussuge.fintech ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2149))

### BUG FIXES (1 change)

- resource/gitlab_project:  [Only apply `use_custom_template` during creation if the value is `true`. Fixes an issue where using `use_custom_template = false` would cause a 500 error](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/bf658535994f38b8dcf743e2827abe4e3e66e20b) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2140))

## 17.4.0 (2024-09-19)

This release was tested against GitLab 17.4, 17.3, and 17.2 for both CE and EE

### FEATURES (4 changes)

- **New Resource** resource/gitlab_member_role: [Allows creating custom roles at an instance or group level](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/5c4980f0f8320c713ae753d1623fedf1c0499ce9) by @slaup1 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2115))
- **New Resource** resource/user_impersonation_token: [Adds support for generating an impersonation token with an admin user](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b0008c5ed590b7f6cc1a52928a75a40d3da6c5b2) by @slaup1 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2086))
- **New Datasource** data/gitlab_protected_tag(s): [Two new datasources allow retrieving a single tag or list of tags by project](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/759279ab7d332e827c87588aa708c38019551e6a) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2126))
- **New Datasource** data/gitlab_pipeline_schedule(s): [Two new datasources allow retrieving a single pipeline schedule or list of pipeline scheduled by project](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/65e4bdb59fd877d1bbe796ce2900d335f60baff0) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2075))

### IMPROVEMENTS (12 changes)

- resource/gitlab_group: [Add support for `default_branch_protection_defaults`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d3a97d371b04c67404cd58d9200a4907a0199cbf) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2113))
- resource/gitlab_group: [Add support for `reject_non_dco_commits`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b9bc021dad13a05d802c8a3c56cbb0b57e5cc5e6) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2105))
- resource/gitlab_project_approval_rule: [Add support for `report_type`, allowing the creation of `Coverage-Check` rules](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/6011774f1ba63c365e7d6dfb21b360e6d5467d08) by @MDeVilliers1 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1913))
- resource/gitlab_group_label: [Update format of the resource ID to allow `name` updates](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/62a853fb14eb6ebea7f31a7504997e8a2068ba47) by @nvh04121 ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2094))
- resource/gitlab_group_label: [Updating `name` no longer forces destroy/creation of the label](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/94f5f9e9d3d389b1d6380ebf659022192a63a45a) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2127)
- resource/gitlab_service_jira: [Aligned fields to upstream API, adding support for 9 new fields](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/f63965d3b5af927f253bb00a29be44e25d6d382c) by @deepflame ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1961) 
- resource/gitlab_application_settings: [Add support for 25 new fields](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/233c3888688b3885d1e716830efbe7d6b4daf018) by @Jitsusama ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2108))
- resource/gitlab_project: [Add support for `pre_receive_secret_detection_enabled`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/38717b856f75908646c8b0a6d0484c04929f04a9) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2106))
- resource/gitlab_project: [Add support for `reject_non_dco_commits`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b9bc021dad13a05d802c8a3c56cbb0b57e5cc5e6) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2105))
- resource/gitlab_project_push_rules: [Add support for `reject_non_dco_commits`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/b9bc021dad13a05d802c8a3c56cbb0b57e5cc5e6) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2105))
- resource/gitlab_user: [Add support for `extern_uid` and `external_provider`](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/bc13a7dd9052e18e8c14beae56eece2dcec13b80) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2097))
- docs/use-case-tech-lead-bootstrapping-small-team: [Added new documentation to the provider with an example of creating a small team using GitLab](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/f90a66f2d288e6556b6174a0f62b35cfb8e2a764) by @Jitsusama ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2112))


### BUG FIXES (2 changes)

- resource/gitlab_user_runner: [Fix an issue where managing an imported runner would fail with an error that `token` was unknown](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/e126dbb7307e2b907b6fd62932d72b33ad788b81) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2114))
- resource/gitlab_personal_access_token: [Fix an issue where tokens with no expiration date encountered a provider error](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/142ce83b0eea8a829c6c7ef04b1bd6b687c0a616) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2104))


## 17.3.1 (2024-08-27)

This release was tested against GitLab 17.3, 17.2, and 17.1 for both CE and EE

### IMPROVEMENTS (1 change)

- [The provider now uses an exponential backoff when no RateLimit-Reset headers are sent from GitLab]() by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2073))

### BUG FIXES (7 changes)

- resource/gitlab_project_access_token: [Fixed an error when applying an update using a plan file that could cause the ID to return an error when rotating](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d0a63fec70fd01d447091bf81138f76692b59c58) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2060))
- resource/gitlab_project_access_token: [Fixed an error when applying an update using a plan from from a previous day to rotate an access token](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d0a63fec70fd01d447091bf81138f76692b59c58) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2060))
- resource/gitlab_group_access_token: [Fixed an error when applying an update using a plan file that could cause the ID to return an error when rotating](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d0a63fec70fd01d447091bf81138f76692b59c58) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2060))
- resource/gitlab_group_access_token: [Fixed an error when applying an update using a plan from from a previous day to rotate an access token](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/d0a63fec70fd01d447091bf81138f76692b59c58) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2060))
- resource/gitlab_group_protected_environment: [Updated documentation for `deploy_access_levels` to be more specific](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/0d01e91b32fadd6980e666b11c877fe97c2ee45a) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2072))
- resource/gitlab_project_protected_environment: [Updated documentation for `deploy_access_levels` to be more specific](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/0d01e91b32fadd6980e666b11c877fe97c2ee45a) by @mness ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2072))
- resource/gitlab_group_saml_link: [Updated documentation for `access_level` to be more descriptive when using custom roles](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/commit/98d0762153fd4ec01d7df62ce5f991bc5b342d82) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2071))

## 17.3.0 (2024-08-15)

This release was tested against GitLab 17.3, 17.2, and 17.1 for both CE and EE

### IMPROVEMENTS (8 changes)

- **New Resource** resource/gitlab_group_security_policy_attachment: [Add new resource for association security policy projects to a group](gitlab-org/terraform-provider-gitlab@fa00e1ec7636dc47b9b4292d164825bf39af0cef) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2048))
- **New Datasource** data/gitlab_group_ids: [Add new datasource retrieving group ID, Full Path, or GraphQL ID based on group ID or Url-encoded path](gitlab-org/terraform-provider-gitlab@e0b0f88070caea6592c1983aa7792a77b27e9b88) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2057))
- **New Datasource** data/gitlab_project_ids: [Add new datasource retrieving project ID, Full Path, or GraphQL ID based on project ID or Url-encoded path](gitlab-org/terraform-provider-gitlab@e0b0f88070caea6592c1983aa7792a77b27e9b88) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2057))
- resource/gitlab_project_job_token_scopes: [Add support for `target_group_ids`, allowing groups to be used in the allow list](gitlab-org/terraform-provider-gitlab@f4246b814952900a5ddfb509ca2c081d47e820e1) by @chapsuK ([merge request](gitlab-org/terraform-provider-gitlab!2040))
- resource/gitlab_group: [Add support for `permanently_remove_on_delete`, allowing groups to be permanently deleted on destroy instead of entering a waiting period](gitlab-org/terraform-provider-gitlab@a0374ef9b9ff87cc9f93a1523fc3bd0998ab0c6f) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2044))
- resource/gitlab_application_settings: [Add support for `max_terraform_state_size_bytes`, allowing configuration of terraform state size at an instance level](gitlab-org/terraform-provider-gitlab@d852a3360f034648816969cb3f94e310f6fb74d9) by @christian.affolter ([merge request](gitlab-org/terraform-provider-gitlab!1960))
- resource/gitlab_project_push_rules: [`gitlab_project_push_rules` will now take over ownership of existing push rules, since a default set is created on project creation](gitlab-org/terraform-provider-gitlab@89e32696080ab4647d49cd077f7ffe0a5675766e) by @mness ([merge request](gitlab-org/terraform-provider-gitlab!2046))
- resource/gitlab_project_job_token_scopes: [Added support for `project`, allowing the use of the resource by project ID or url-encoded path](gitlab-org/terraform-provider-gitlab@f040983bed1545fe445d41fb56e36f938cee35de) by @heidi.berry ([merge request](gitlab-org/terraform-provider-gitlab!1981))

### BUG FIXES (4 changes)

- resource/gitlab_group_access_token: [Fix example to use the proper resource](gitlab-org/terraform-provider-gitlab@f20836445c6190feb4dc081eca9f02887b39c205) by @blrz ([merge request](gitlab-org/terraform-provider-gitlab!2050)) 
- resource/gitlab_group_security_policy_attachment: [Fix example to use the proper resource](gitlab-org/terraform-provider-gitlab@0fe493c436fb7bb333d40c8f5aab8da18f1c272c) by @calebw ([merge request](gitlab-org/terraform-provider-gitlab!2054)) 
- resource/gitlab_project_security_policy_attachment: [Fix example to use the proper resource](gitlab-org/terraform-provider-gitlab@0fe493c436fb7bb333d40c8f5aab8da18f1c272c) by @calebw ([merge request](gitlab-org/terraform-provider-gitlab!2054)) 
- resource/gitlab_user_runner: [Fix description to remove typo in `group_type`](gitlab-org/terraform-provider-gitlab@f67c731da1f9fcd102053c047de1455043853291) by @clemenschuaccso ([merge request](gitlab-org/terraform-provider-gitlab!2049)) 

## 17.2.0 (2024-07-18)

This release was tested against GitLab 17.2, 17.1, and 17.0 for both CE and EE

### BREAKING CHANGES (1 change)
Note: As a security related change, this breaking change is allowed outside a major release. If a non-sensitive token is required for migration purposes users can use the `nonsensitive()` TF function.

- resource/gitlab_cluster_agent_token: [Mark GitLab cluster agent token as sensitive](gitlab-org/terraform-provider-gitlab@1eec1065723f6393a3e5deb45356a47ec27b575e) by @CarbonCollins ([merge request](gitlab-org/terraform-provider-gitlab!2032)) 

### IMPROVEMENTS (1 change)

- **New Resource** resource/gitlab_project_security_policy_attachment: [Add new resource for associating security policy projects to a project](gitlab-org/terraform-provider-gitlab@3fc5f6c8a9d74f8b034af84a73fb9dcb6d323179) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2009))
- resource/gitlab_personal_access_token: [Migrate `gitlab_personal_access_token` from SDK to TF Plugin Framework](gitlab-org/terraform-provider-gitlab@24489118da1ce464994cdfcedecbe804edf65490) by @theipster ([merge request](gitlab-org/terraform-provider-gitlab!2005)) 

### BUG FIXES (4 changes)

- resource/gitlab_deploy_key: [Fix deploy key documentation reference to point to the correct resource for enabling a pre-existing deploy key](gitlab-org/terraform-provider-gitlab@f8718e66a5608c521009d091761b5577a055e904) by @blrz ([merge request](gitlab-org/terraform-provider-gitlab!2024)) 
- resource/gitlab_user_runner: [Fix an issue with gitlab_user_runner causing inconsistent TF plans](gitlab-org/terraform-provider-gitlab@a5448baf255134f0d0bff5fe16803ee4cf35ce61) by @PatrickRice ([merge request](gitlab-org/terraform-provider-gitlab!2020)) 

## 17.1.0 (2024-06-20)

This release was tested against GitLab 17.1, 17.0, and 16.11 for both CE and EE

### BUG FIXES (3 changes)

- resource/gitlab_project_job_token_scopes:  [Correct examples for resource project_job_token_scopes](gitlab-org/terraform-provider-gitlab@3027ed1e0861b405ef81f13b3e58fce5aff96ed2) by @heidi.berry ([merge request](gitlab-org/terraform-provider-gitlab!1968)) 
- resource/gitlab_group_access_token: [Fix an issue where using access tokens with a `time` provider caused segfaults](gitlab-org/terraform-provider-gitlab@de3b98c0881a8aa8a57ca9ea2ac735230b56f8c2) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1999))
- resource/gitlab_project_access_token: [Fix an issue where using access tokens with a `time` provider caused segfaults](gitlab-org/terraform-provider-gitlab@de3b98c0881a8aa8a57ca9ea2ac735230b56f8c2) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1999))

### IMPROVEMENTS (8 changes)

- **New Resource** resource/gitlab_telegram_integration: [Implement Telegram Integration](gitlab-org/terraform-provider-gitlab@bc5eb7cf0138296aedb6f2559166a95f3f56bf32) by @alxrem ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1963))
- resource/gitlab_group_saml_link: [Add ability to set custom roles on group SAML links](gitlab-org/terraform-provider-gitlab1cd431c8238fc690e04e9625951208693977e2c0) by @heidi.berry ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1949))
- resource/gitlab_integration_slack:  [Add confidential_note_channel to slack integration](gitlab-org/terraform-provider-gitlab@38c49916a0d21151a59de518ca103b1fca3c5d6d) by @0oMarko0 ([merge request](gitlab-org/terraform-provider-gitlab!1988)) 
- resource/gitlab_personal_access_token:  [Implement manage_runner access token scope](gitlab-org/terraform-provider-gitlab@a5a1c711d0d086a51dfd92fcc751814182e526b8) ([merge request](gitlab-org/terraform-provider-gitlab!1986)) 
- resource/gitlab_personal_access_token:  [Mark `expires_at` attribute as optional](gitlab-org/terraform-provider-gitlab@7be1aa7f6948552eb2e16331b80c83e3f26eb971) by @erezo9 ([merge request](gitlab-org/terraform-provider-gitlab!1983))
- resource/gitlab_group_access_token: [Fixed several documentation issues with access token resources](gitlab-org/terraform-provider-gitlab@738a1ffb41ad7cfaffa205d0c31f827d9f4e59e8) by @theipster ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2004))
- resource/gitlab_project_access_token: [Fixed several documentation issues with access token resources](gitlab-org/terraform-provider-gitlab@738a1ffb41ad7cfaffa205d0c31f827d9f4e59e8) by @theipster ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/2004))
- resource/gitlab_project: [Update `initialize_with_readme` to add border case documentation](gitlab-org/terraform-provider-gitlab@247453b3e45362c4ef91fa9f537d6de7cf49cbfc) by @PatrickRice ([merge request](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1989))


## 17.0.1 (2024-06-07)

This release was tested against GitLab 17.0, 16.11, and 16.10 for both CE and EE

BUG FIXES:

- Fixed an issue where sensitive tokens were not masked properly in debug log files ([!1997](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1997))

## 17.0.0 (2024-05-16)

This release was tested against GitLab 17.0, 16.11, and 16.10 for both CE and EE

KNOWN ISSUES:

- `gitlab_current_user` returns an empty string for `public_email` ([#6305](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/issues/6305))

BREAKING CHANGES:

- resources/project_protected_environment: Removed support for `required_approval_count` field, use `required_approvals` in `approval_rules` or `deploy_access_level` instead ([!1940](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1940))
- resources/group_protected_environment: Removed support for `required_approval_count` field, use `required_approvals` in `approval_rules` or `deploy_access_level` instead ([!1940](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1940))
- resources/gitlab_group: Removed a version check related to `commit_committer_check` and `reject_unsigned_commits` that would prevent their use in versions prior to GitLab 16.4. If used with versions earlier than 16.4, these attributes will cause an error instead of being excluded. ([!1937](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1937))
- resources/gitlab_group: Removed support for `emails_disabled`, use `emails_enabled` instead ([!1929](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1929))
- resources/gitlab_project: Removed support for `emails_disabled`, use `emails_enabled` instead ([!1929](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1929))
- resources/gitlab_pipeline_schedule: `ref` now requires the full ref instead of allowing the use of the short ref. If you previously used `main`, you now need to use `refs/heads/main` instead, for example ([!1923](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1923))
- resources/gitlab_pipeline_trigger: `token` can no longer be imported. ([!1905](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1905))
- resources/gitlab_pipeline_trigger: Updating the `project` attribute will now force the creation of a new pipeline trigger ([!1905](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1905))
- data/gitlab_group(s): Removed support for `emails_disabled`, use `emails_enabled` instead ([!1929](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1929))
- data/gitlab_project(s): Removed support for `emails_disabled`, use `emails_enabled` instead ([!1929](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1929))
- data/gitlab_project(s): Removed support for `public`, use `visibility` instead ([!1929](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1929))
- `master` is no longer a valid access level on any resource that supports the use of access levels. This impacts the resources listed below. ([!1903](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1903))
  - gitlab_group_access_token
  - gitlab_group_ldap_link
  - gitlab_group_membership
  - gitlab_group_share_group
  - gitlab_project_access_token
  - gitlab_project_membership
  - gitlab_project_share_group

IMPROVEMENTS:

- **New Resource** resource/gitlab_integration_jenkins: Allows managing a project Jenkins integration ([!1919](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1919))
- **New Resource** resource/gitlab_project_push_rules:  Allows managing the lifecycle of push rules on a project ([!1893](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1893))
- **New Resource** resource/gitlab_project_job_token_scopes: This resource sets a strict list of project job token scopes, and removes any job token scopes not managed by the resource. This can be useful to explicitly deny job token scopes on a project. ([!1907](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1907))
- resources/gitlab_pipeline_schedule_variable: Added support for the use of `variable_type` ([!1952](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1952))
- resources/gitlab_group: Added support for the use of `commit_committer_name_check` to the `push_rules` block ([!1937](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1937))
- resources/gitlab_project: Added support for the use of `commit_committer_name_check` to the `push_rules` block ([!1918](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1918))
- resources/gitlab_instance_variable: Added support for the use of `description` ([!1950](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1950))
- resources/gitlab_user_runner: Added example documentation for this resource to make it easier to consume ([!1928](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1928))
- resources/gitlab_application_settings: Add support for `minimum_password_length` to the resource ([!1917](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1917))
- resources/gitlab_personal_access_token: Updated the API used to read personal access token data, which improves performance of this resource in situations where many tokens are being maintained, and improves reliability of the resource in high concurrency situations ([!1908](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1908))
- data/gitlab_instance_variable: Added support for `description` ([!1950](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1950))
- provider: Added documentation that the use of Project Access Tokens or Group Access Tokens may not work with all resources ([!1928](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1928))

BUG FIXES:

- resources/gitlab_pipeline_schedule: Fixed a provider crash in situations where the scheduled pipeline fails to create ([!1899](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1899))
- resources/gitlab_group: Fixed a provider error when attempting to create groups with `push_rules` on GitLab CE where `push_rules` are not supported ([!1891](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1891))
- resources/gitlab_group_access_token: Fixed an issue with token rotation using `rotation_configuration` where tokens wouldn't rotate properly after `expires_at` was stored in state. Added additional debug logging for token rotation. ([!1953](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1953))
- resources/gitlab_group_access_token: Fixed an issue with token rotation where manually managing expiration using `expires_at` would encounter an error after updating `expires_at` twice ([!1916](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1916))
- resources/gitlab_project_access_token: Fixed an issue with token rotation using `rotation_configuration` where tokens wouldn't rotate properly after `expires_at` was stored in state. Added additional debug logging for token rotation. ([!1953](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1953))
- resources/gitlab_project_access_token: Fixed an issue with token rotation where manually managing expiration using `expires_at` would encounter an error after updating `expires_at` twice ([!1916](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1916))
- resources/gitlab_pipeline_trigger: Fixed a bug where applying TF with different users could corrupt the pipeline trigger `token` [!1905](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1905)

## 16.11.0 (2024-04-18)

This release was tested against GitLab 16.9, 16.10, and 16.11 for both CE and EE

IMPROVEMENTS:

- **New Data Source** data/gitlab_compliance_framework: Allows querying Compliance Frameworks to help retrieve the ID for use in downstream resources ([!1880](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1880))
- resources/gitlab_project_access_token: Added support for the use of `rotation_configuration` to automatically rotate tokens periodically. ([!1887](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1887))
- resources/gitlab_project_access_token: Added support for rotating the token by changing the `expires_at` instead of deleting and re-creating the token. ([!1887](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1887))
- resources/gitlab_group_access_token: Added support for the use of `rotation_configuration` to automatically rotate tokens periodically. ([!1887](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1887))
- resources/gitlab_group_access_token: Added support for rotating the token by changing the `expires_at` instead of deleting and re-creating the token. ([!1887](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1887))
- resources/gitlab_project_access_token: Added support for new token scopes related to AI, k8s, and observability ([!1878](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1878))
- resources/gitlab_group_access_token: Added support for new token scopes related to AI, k8s, and observability ([!1878](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1878))
- resources/gitlab_project: Added support for `emails_enabled` and deprecated support for `emails_disabled`, which will be removed in 17.0 ([!1881](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1881))
- resources/gitlab_project_protected_environment: Added support for `group_inheritance_type` ([!1855](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1855))
- resources/gitlab_group_protected_environment: Added support for `group_inheritance_type` ([!1855](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1855))
- resources/gitlab_project_hook: Added support for `custom_webhook_template` ([!1862](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1862))
- resources/gitlab_group_hook: Added support for `custom_webhook_template` ([!1862](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1862))
- resources/gitlab_group_membership: Added support for `member_role_id`, enabling the use of a custom role when assigning users to a group ([!1809](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1809))
- data/gitlab_project_hook(s): Added support for `custom_webhook_template` ([!1862](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1862))
- data/gitlab_group_hook(s): Added support for `custom_webhook_template` ([!1862](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1862))

BUG FIXES:

- resource/gitlab_project_hook: Fixed an issue where changing the `project` value didn't force a new resource ([!1871](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1871))

## 16.10.0 (2024-03-21)

This release was tested against GitLab 16.8, 16.9, and 16.10 for both CE and EE

NOTES:

- scripts/gitlab.rb has been updated for local development to set the license mode and customer portal URL for testing. If you're using a personal license for local development, you may need to update this file temporarily to run EE locally. ([!1861](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1861))

IMPROVEMENTS:
- **New Data Source** datasource/gitlab_release: Allows querying a GitLab Release by project and tag name to get release information or assets ([!1851](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1851))

BUG FIXES:

- resources/gitlab_project_level_mr_approvals: fixed a documentation issue with the use of `merge_requests_disable_committers_approval` ([!1864](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1864))
- resources/gitlab_group_ldap_link: fixed an issue where deleting the group associated to an LDAP link would result in a TF state that required manual intervention. Using a value of `true` with the `force` attribute will now remove the LDAP link from state when the group is deleted. ([!1842](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1842))
- resources/gitlab_cluster_agent: fixed the example documentation to show the file contents as encoded ([!1852](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1852))
- resources/gitlab_pipeline_schedule: fixed a potential panic on the provider that could occur when there was an error editing the pipeline schedule ([!1847](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1847))


## 16.9.1 (2024-02-15)

This release was tested against GitLab 16.7, 16.8, and 16.9 for both CE and EE

BUG FIXES:

- Fixed a Go version mismatch when using `goreleaser` that prevented v16.9.0 from being published properly ([!1839](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1839))

## 16.9.0 (2024-02-15)

This release was tested against GitLab 16.7, 16.8, and 16.9 for both CE and EE

IMPROVEMENTS:
- resource/gitlab_project_variable: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- resource/gitlab_group_variable: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- resource/gitlab_project: added support for `ci_restrict_pipeline_cancellation_role` ([!1825](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1825))
- datasource/gitlab_project_variable: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- datasource/gitlab_project_variables: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- datasource/gitlab_group_variable: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- datasource/gitlab_group_variables: added support for `description` ([!1827](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1827))
- datasource/gitlab_project: added support for `ci_restrict_pipeline_cancellation_role` ([!1825](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1825))
- datasource/gitlab_projects: added support for `ci_restrict_pipeline_cancellation_role` ([!1825](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1825))

BUG FIXES:

- resource/gitlab_project_protected_environment: Fixed an issue where using characters in the project name that required encoding (such as "/") would cause an error ([!1835](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1835))
- resource/application_settings: Fixed an issue where elasticsearch namespace and project IDs accepted a list of strings instead of a list of integers. Providing a list of integers would cause a provider error, rendering the attributes unusable. ([!1824](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1824))

## 16.8.1 (2024-01-24)

This release was tested against GitLab 16.6, 16.7, and 16.8 for both CE and EE

BUG FIXES:

- resource/gitlab_compliance_framework: Improved permissions-based error handling for the resource, so permissions errors don't result in a "provider error" warning ([!6193](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/issues/6193))
- resource/gitlab_pipeline_schedule: Fixed an issue with `take_ownership = true` where the ownership would only be updated when a separate attribute change was identified. The provider will now always assume ownership even if no other changes are identified ([!1765](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1765))

## 16.8.0 (2024-01-18)

This release was tested against GitLab 16.6, 16.7, and 16.8 for both CE and EE

BREAKING CHANGE:

- gitlab_application_settings: Removed support for `delayed_group_deletion` and `delayed_project_deletion`, which haven't been supported since GitLab 16.0 ([!1799](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1799))

IMPROVEMENTS:

- **New Resource** `gitlab_global_level_notifications` allows managing global notifications for the user ([!1801](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1801))
- resource/gitlab_project: Added support for the `timeouts` block, to allow configurable timeouts for creating projects ([!1797](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1797))
- resource/gitlab_application_settings: Added support for `housekeeping_optimize_repository_period`, and updated the description of `housekeeping_enabled` to be more descriptive of what fields were required for a successful apply ([!1777](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1777))
- datasource/gitlab_project: Added support for `shared_with_groups` which includes groups that are shared with the project ([!1795](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1795))
- datasource/gitlab_group: Added support for `shared_with_groups` which includes groups that are shared with the group ([!1769](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1769))
- datasource/gitlab_user: Added support for the `is_bot` attribute ([!1798](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1798))
- datasource/gitlab_users: Added support for the `is_bot` attribute ([!1798](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1798))

BUG FIXES:

- resource/gitlab_pipeline_schedule_variable: Fixed an issue where deleting the pipeline associated to the variable outside of terraform would cause the resource to be stuck in an error state ([!1796](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1796))
- resource/gitlab_project_hook: Fixed an issue where changing the URL of a webhook would delete the associated `token`. Changing the URL will now force a new webhook to be created ([!1794](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1794))

## 16.7.0 (2024-01-08)

This release was tested against GitLab 16.5, 16.6, and 16.7 for both CE and EE

DEPRECATION:

- This release updates the recommended Terraform version for this Provider from 1.0.0 to 1.4.0, and updates the version of Terraform we use for CI/CD to 1.4.0 as a result. This is related to a bug that was fixed in Terraform 1.4.0 related to how complex objects are compared. Without using Terraform 1.4.0, the provider cannot guarantee that plan output using nested objects is the same every time. Prior versions will likely still result in a successul plan and apply, but we will ask you to update prior to assisting with issue triage.

IMPROVEMENTS:

- resource/gitlab_branch_protection: Updating `allowed_to_push` will no longer destroy and re-create branch protection, it will instead update it in-place ([!1593](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1593))
- resource/gitlab_group_issue_board: Issue boards now supports the use of scoped labels, and label position can be explicitly configured ([!1771](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1771))
- resource/gitlab_project: Add support for `group_runners_enabled` ([!1735](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1735))
- datasource/gitlab_group: Groups with many projects will now be retrieved significantly faster ([!1770](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1770))
- datasource/gitlab_project: Add support for `group_runners_enabled` ([!1735](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1735))

BUG FIXES:

- resource/gitlab_pipeline_schedule: Fixed an issue where a pipeline schedule with no owner could cause a provider panic ([!1762](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1762))
- resource/gitlab_group_ldap_link: Fixed an issue with `force` that could cause an error when attempting to delete an ldap link ([!1757](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1757))

## 16.6.0 (2023-11-16)

This release was tested against GitLab 16.4, 16.5, and 16.6 for both CE and EE

KNOWN ISSUES:

- Attempting to use the `gitlab_users` datasource with `sort` will not return users in the specified sort order when used with GitLab 16.6.0, as GitLab 16.6.0 uses relevancy sorting and ignores `sort`. This will be resolved with GitLab 16.6.1.

IMPROVEMENTS:

- **New Resource:** `gitlab_project_level_notifications` allows managing notification events for project ([!1715](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1715))
- resource/gitlab_project_approval_rule: added support for `applies_to_all_protected_branches` ([!1755](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1755))
- resource/gitlab_pipeline_schedule: added support for `take_ownership`, which will take ownership of the pipeline schedule prior to attempting an update ([!1745](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1745))
- resource/gitlab_group: added support for `push_rules` ([!1730](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1730))

BUG FIXES:

- resource/gitlab_user_runner: Fixed an issue where not including `maximum_timeout` could cause an issue when updating the runner ([!1758](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1758))
- datasource/gitlab_user: When using `email`, the the data source will now return the first user returned from the API instead of encountering an error when more than one is identified. When used with GitLab 16.6.0, this will always be the exact match if an exact match is available. ([!1743](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1743))


## 16.5.0 (2023-10-22)

This release was tested against GitLab 16.3, 16.4, and 16.5 for both CE and EE

IMPROVEMENTS:

- **New Resource:** `gitlab_group_protected_environment` allows managing group-level protected environments ([!1707](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1707))
- resource/gitlab_user_sshkey: Added support for creating an SSH key for the current user by making `user_id` optional ([!1726](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1726))
- resource/gitlab_group: Added support for managing the `shared_runners_setting` attribute ([!1710](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1710))
- resource/gitlab_project: Added support for creating an empty repository using the `empty_repo` attribute ([!1713](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1713))
- resource/gitlab_project: Added support for the `public_jobs` attribute, deprecating the old `public_builds` attribute ([!1700](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1700))
- datasource/gitlab_project: Added support for reading the `empty_repo` attribute ([!1713](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1713))
- datasource/gitlab_projects: Added support for reading the `empty_repo` attribute ([!1713](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1713))
- datasource/gitlab_group: Added support for reading the `shared_runner_setting` attribute ([!1717](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1717))
- datasource/gitlab_groups: Added support for reading the `shared_runner_setting` attribute ([!1717](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1717))
- datasource/gitlab_group_subgroups: Added support for reading the `shared_runner_setting` attribute ([!1719](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1719))

BUG FIXES:

- resource/gitlab_group: Removed "default" hints in the documentation, since defaults can be changed by admins in some cases ([!1696](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1696))
- resource/gitlab_group_ldap_link: Fixed an issue where changing CN or Filter didn't force a new resource ([!1729](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1729))
- resource/gitlab_project: Fixed an issue where the documentation didn't contain valid values for several fields ([!1714](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1714))
- resource/gitlab_tag_protection: Fix d an issue where the resource read the wrong `create_access_level` when using `no one` ([!1694](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1694))

## 16.4.1 (2023-09-25)

This release was tested against GitLab 16.2, 16.3, and 16.4 for both CE and EE

BUG FIXES:

- resource/gitlab_project_protected_environment: Fix segfault when using `for_each` for `deploy_access_levels` when `approval_rules` are not specified ([!1699](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1699))

## 16.4.0 (2023-09-22)

This release was tested against GitLab 16.2, 16.3, and 16.4 for both CE and EE

BREAKING CHANGES:

This breaking change was made early for security reasons. If a configuration relies on the value being non-sensitive,
users can use the [`nonsensitive()`](https://developer.hashicorp.com/terraform/language/functions/nonsensitive) function
in Terraform.

- resource/gitlab_user_runner: `token` is now marked as sensitive ([!1688](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1688))

IMPROVEMENTS:

- resource/gitlab_project_mirror: Updated documentation to include a warning about `keep_divergent_refs` default value ([!1691](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1691))
- resource/gitlab_project_protected_environment: Add support for `approval_rules` ([!1679](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1679))
- resource/gitlab_group_access_token: Add support for the `create_runner` scope ([!1675](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1675))
- resource/gitlab_personal_access_token: Add support for the `create_runner` scope ([!1675](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1675))
- resource/gitlab_project_access_token: Add support for the `create_runner` scope ([!1675](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1675))

BUG FIXES:

- resource/gitlab_branch: No longer returns an error when the branch is missing during a destroy ([!1690](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1690))
- datasource/gitlab_cluster_agents: Fixed an issue where `agent_id` was always `0` ([!1677](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1677))
- datasource/gitlab_group_subgroups: Fixed an issue where the data source returned a maximum of 20 subgroups ([!1689](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1689))

## 16.3.0 (2023-08-22)

This release was tested against GitLab 16.0, 16.1, and 16.2 for both CE and EE

IMPROVEMENTS:

- **New Resource:** `gitlab_user_runner` allows managing runners using the new runner flow without using a registration token ([!1618](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1618))
- **New Resource:** `gitlab_group_epic_board` allows managing epic boards for groups ([!1658](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1658))
- **New Resource:** `gitlab_project_job_token_allow` allows managing the inbound allow list for a project when using Job Tokens ([!1631](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1631))
- resource/repository_file: Add the ability to specify a different commit message for Create/Update/Delete operation ([!1629](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1629))
- resource/gitlab_project_level_mr_approvals: Add support for `selective_code_owner_removals` ([!1641](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1641))
- resource/gitlab_group: Add support for `wiki_access_level` ([!1656](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1656))
- resource/gitlab_group_badge: Add support for `name` ([!1655](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1655))
- datasource/gitlab_group: Add support for `wiki_access_level` ([!1656](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1656))
- datasource/gitlab_groups: Add support for `wiki_access_level` ([!1656](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1656))


BUG FIXES:

- resource/gitlab_group_access_token: Require the `expires_at` attribute ([!1661](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1661))
- resource/gitlab_personal_access_token: Require the `expires_at` attribute ([!1661](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1661))
- resource/gitlab_project_access_token: Require the `expires_at` attribute ([!1661](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1661))
- resource/gitlab_pipeline_schedule_variable: Fix several spelling errors in the documentation that would make examples non-functional ([!1647](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1647))


## 16.2.0 (2023-07-22)

This release was tested against GitLab 15.11, 16.0 and 16.1 for both CE and EE.

IMPROVEMENTS:

- **New Resource:** `gitlab_project_compliance_framework` ([!1616](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1616))
- **New Resource:** `gitlab_compliance_framework` ([!1599](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1599))
- resource/gitlab_group_badge: Improve examples to contain common badges ([!1627](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1627))
- resource/branch_protection: Add support for `admin` as a value for `unprotect_access_level` ([!1626](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1626))
- datasource/gitlab_groups: Add `top_level_only` support ([!1606](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1606))
- datasource/gitlab_project: Add `topic` support ([!1610](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1610))

BUG FIXES:

- resource/gitlab_application_settings: Fix documentation to list the correct `import_sources` ([!1638](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1638))
- resource/gitlab_project: Previously, `name_regex_delete` was improperly deprecated. Remove deprecation notice, and add notice to `name_regex`, which is the proper field ([!1600](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1600))
- resource/gitlab_repository_file: Fix an issue where updating a repository file when using the `text` encoding returned a base64 encoding error ([!1642](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1642))

## 16.1.1 (2023-07-17)

This release was tested against GitLab 15.11, 16.0 and 16.1 for both CE and EE.

IMRPOVEMENTS:

- resource/gitlab_repository_file: Support the use of the `encoding` parameter. This allows using `text` encoding, which re-introduces the ability to have plaintext comparisons during a plan operation ([!1633](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1633))

## 16.1.0 (2023-06-22)

This release was tested against GitLab 15.11, 16.0 and 16.1 for both CE and EE.

IMRPOVEMENTS:

- resource/gitlab_tag_protection: Support `allowed_to_create` attribute ([!1549](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1549))

BUG FIXES:

- resource/gitlab_deploy_token: Fix dynamic `username` attribute reading after creation ([!1569](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1569))
- resource/gitlab_branch_protection: Remove unsupported `no one` value for `unprotect_access_level` ([!1594](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1594))
- resource/gitlab_project: Fix deprecation of `name_regex` attribute ([!1600](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1600))

MISC:

- Added a new guide for how to contribute a new resource from scratch ([!1487](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1487))

## 16.0.3 (2023-05-24)

This release was tested against GitLab 15.10, 15.11 and 16.0 for both CE and EE.

BREAKING CHANGES:

Since this is a bug fix release for a major release with breaking changes you may
want to follow the [Terraform GitLab Provider Version 16.0 Upgrade Guide](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs/guides/version-16.0-upgrade) for details.

BUG FIXES:

- resource/gitlab_group_ldap_link: Fix state migration for `group_id` to `group`

## 16.0.2 (2023-05-23)

This release was tested against GitLab 15.10, 15.11 and 16.0 for both CE and EE.

BREAKING CHANGES:

Since this is a bug fix release for a major release with breaking changes you may
want to follow the [Terraform GitLab Provider Version 16.0 Upgrade Guide](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs/guides/version-16.0-upgrade) for details.

BUG FIXES:

- resource/gitlab_project_level_mr_approvals: Fix state migration for `project_id` to `project` when already on v16
- resource/gitlab_project_freeze_period: Fix state migration for `project_id` to `project` when already on v16
- resource/gitlab_project_membership: Fix state migration for `project_id` to `project` when already on v16
- resource/gitlab_project_share_group: Fix state migration for `project_id` to `project` when already on v16

## 16.0.1 (2023-05-23)

This release was tested against GitLab 15.10, 15.11 and 16.0 for both CE and EE.

BREAKING CHANGES:

Since this is a bug fix release for a major release with breaking changes you may
want to follow the [Terraform GitLab Provider Version 16.0 Upgrade Guide](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs/guides/version-16.0-upgrade) for details.

BUG FIXES:

- resource/gitlab_pipeline_schedule_variable: Fix panic when `pipeline_schedule_id` is a `float64` in state
- resource/gitlab_project_variable: Fix panic when upgrading from a state that was created prior to GitLab 13.4
- resource/gitlab_project_level_mr_approvals: Add state migration for `project_id` to `project`
- resource/gitlab_project_freeze_period: Add state migration for `project_id` to `project`
- resource/gitlab_project_membership: Add state migration for `project_id` to `project`
- resource/gitlab_project_share_group: Add state migration for `project_id` to `project`
- resource/gitlab_project_access_token: Mark `expires_at` as computed
- resource/gitlab_group_access_token: Mark `expires_at` as computed

## 16.0.0 (2023-05-22)

- This release was tested against GitLab 15.10, 15.11 and 16.0 for both CE and EE.
- **Note:** this is a major release and breaks some interfaces in resources and data sources of
  this provider, but also drops support for older GitLab versions.

BREAKING CHANGES:

See [Terraform GitLab Provider Version 16.0 Upgrade Guide](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs/guides/version-16.0-upgrade) for details.

BREAKING CHANGES:

- resource/gitlab_instance_variable: Change `value` attribute to non-sensitive ([!1521](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1521))
- resource/gitlab_group_variable: Change `value` attribute to non-sensitive ([!1521](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1521))
- resource/gitlab_project_variable: Change `value` attribute to non-sensitive ([!1521](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1521))
- resource/gitlab_deploy_token: Change resource id format to `<token-type>:<type-id>:<token-id>` ([!1523](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1523))
- resource/gitlab_deploy_key: Change resource id format to `<project>:<key-id>` ([!1522](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1522))
- resource/gitlab_project_hook: Change resource id format to `<project>:<hook-id>` ([!1483](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1483))
- resource/gitlab_group_label: Change resource id format to `<group>:<label-name>` ([!1525](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1525))
- resource/gitlab_label: Rename resource to `gitlab_project_label` ([!1526](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1526))
- resource/gitlab_project_label: Change resource id format to `<project>:<label-name>` ([!1526](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1526))
- resource/gitlab_managed_license: Remove resource ([!1512](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1512))
- resource/gitlab_pipeline_schedule_variable: Change resource id format to `<project>:<schedule-id>:<variable-id>` ([!1529](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1529))
- resource/gitlab_repository_file: Remove support for auto-encoding logic ([!1530](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1530))
- resource/gitlab_group_ldap_link: Change resource id format to `<group>:<provider>:[cn]:[filter]` ([!1527](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1527))
- resource/gitlab_group_ldap_link: Rename `group_id` attribute to `group` ([!1532](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1532))
- resource/gitlab_service_*: Rename `gitlab_service_*` resources to `gitlab_integration_*` ([!1534](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1534))
- resource/gitlab_project: Remove `operations_access_level` ([!1548](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1548))
- datasource/gitlab_project: Remove `operations_access_level` ([!1548](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1548))
- datasource/gitlab_projects: Remove `operations_access_level` ([!1548](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1548))
- resource/gitlab_pipeline_trigger: Change resource id format to `<project>:<trigger-id>` ([!1551](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1551))
- resource/gitlab_pipeline_schedule: Change resource id format to `<project>:<schedule-id>` ([!1551](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1551))
- resource/gitlab_project_freeze_period: Change `project_id` attribute to `project` ([!1553](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1553))
- resource/gitlab_project_level_mr_approvals: Change `project_id` attribute to `project` ([!1553](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1553))
- resource/gitlab_project_membership: Change `project_id` attribute to `project` ([!1553](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1553))
- resource/gitlab_project_share_group: Change `project_id` attribute to `project` ([!1553](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1553))
- resource/gitlab_project_access_token: Require `expires_at` attribute ([!1557](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1557))

FEATURES:

- resource/gitlab_instance_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))
- resource/gitlab_group_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))
- resource/gitlab_project_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))
- datasource/gitlab_instance_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))
- datasource/gitlab_group_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))
- datasource/gitlab_project_variable: Support `raw` attribute ([!1533](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1533))

IMPROVEMENTS:

- resource/gitlab_project_environment: Wait for environment to stop before deleting it ([!1509](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1509))

## 15.11.0 (2023-04-22)

This release was tested against GitLab 15.9, 15.10 and 15.11 for both CE and EE.

BREAKING CHANGES:

- resource/gitlab_project: Remove specialized branch protection logic for GitLab prior to 14.11 ([!1486](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1486))

IMPROVEMENTS:

- resource/gitlab_application_settings: Support `can_create_group` attribute ([!1484](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1484))
- resource/gitlab_project: Support `keep_latest_artifact` attribute ([!1506](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1506))
- resource/gitlab_group_subgroups: Support `skip_groups` attribute ([!1516](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1516))
- datasource/gitlab_project: Support `keep_latest_artifact` attribute ([!1506](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1506))
- datasource/gitlab_projects: Support `keep_latest_artifact` attribute ([!1506](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1506))

BUG FIXES:

- resource/gitlab_project: Fix waiting when `skip_wait_for_default_branch_protection` is disabled ([!1489](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1489))
- resource/gitlab_group_ldap_link: Fix re-creating LDAP link if it was removed out of bounds ([!1495](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1495))
- resource/gitlab_runner: Fix plan for `tags` attribute to make it order independent ([!1492](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1492))


## 15.10.0 (2023-03-22)

This release was tested against GitLab 15.8, 15.9 and 15.10 for both CE and EE.

FEATURES:

- **New Resource:** `gitlab_application` ([#1436](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1436))
- **New Resource:** `gitlab_service_custom_issue_tracker` ([#1459](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1459))

IMPROVEMENTS:

- resource/gitlab_personal_access_token: Support `admin_mode` as value in the `scopes` attribute ([#1456](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1456))
- resource/gitlab_group: Remove explicit attribute defaults. This solves an issue where newly imported resources may have plan changes for default attributes. In some cases, the new API-based default values are more permissive than the old provider-based defaults. See The linked MR for details. ([#1479](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1479))
- resource/gitlab_repository_file: Add validator to avoid leading `/` and `./` in path ([#1472](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1472))
- resource/gitlab_project: Deprecate the `name_regex_delete` in favor of the `name_regex` attribute ([#1466](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1466))
- resource/gitlab_project: Support `environments_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- resource/gitlab_project: Support `feature_flags_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- resource/gitlab_project: Support `infrastructure_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- resource/gitlab_project: Support `monitor_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- resource/gitlab_project: Support `release_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_project: Mark `runners_token` as sensitive ([#1461](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1461))
- datasource/gitlab_project: Support `environments_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_project: Support `feature_flags_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_project: Support `infrastructure_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_project: Support `monitor_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_project: Support `release_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_projects: Support `environments_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_projects: Support `feature_flags_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_projects: Support `infrastructure_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_projects: Support `monitor_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))
- datasource/gitlab_projects: Support `release_access_level` attribute ([#1469](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1469))

BUG FIXES:

- resource/gitlab_project: Add validator to the `id` attribute to prevent misuse ([#1476](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1476))
- resource/gitlab_group: Wait for the group to be fully created and functional ([#1465](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1465))
- resource/gitlab_service_slack: Fix perpetual diff in `webhook` attribute due to upstream API deprecation ([#1470](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1470))
- resource/gitlab_service_microsoft_teams: Fix perpetual diff in `webhook` attribute due to upstream API deprecation ([#1470](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1470))

## 15.9.0 (2023-02-22)

This release was tested against GitLab 15.7, 15.8 and 15.9 for both CE and EE.

FEATURES:

- **New Resource:** `gitlab_pages_domain` ([#1419](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1419))

IMPROVEMENTS:

- provider: Support `GITLAB_EARLY_AUTH_CHECK` environment variable as default for the `early_auth_check` provider attribute ([#1455](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1455))
- resource/gitlab_project: Support public and private repositories in `import_url` for imports and pull mirrors ([#1452](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1452))

BUG FIXES:

- resource/gitlab_application_settings: Support `"nil"` for `enabled_git_protocols` to disable it ([#1457](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1457))
- resource/gitlab_repository_file: Handle `404`s when `overwrite_on_create` is `true` ([#1433](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1433))
- datasource/gitlab_project: Don't get `push_rules` when user does not have permissions for it ([#1450](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1450))

## 15.8.0 (2023-01-22)

This release was tested against GitLab 15.6, 15.7 and 15.8 for both CE and EE.

FEATURES:

- provider: mask `token` in provider logs ([#1394](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1394))

IMPROVEMENTS:

- resource/gitlab_project: Add `ip_restriction_ranges` attribute ([#1392](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1392))
- datasource/gitlab_group_membership: Add `inherited` attribute to include inherited memberships ([#1402](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1402))

BUG FIXES:

- provider: Fix `early_auth_check` always set to `true` ([#1414](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1414))
- resource/gitlab_project_share_group: Fix refreshing shared groups of a project ([#1412](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1412))
- resource/gitlab_project_protected_environment: Fix perpetual diffs of `deploy_access_levels` ([#1421](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1421))
- resource/gitlab_project_approval_rule: auto-import default `any_approver` rule during create ([#1425](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1425))
- resource/gitlab_group: Mark `shared_runners_minutes_limit` and `extra_shared_runners_minutes_limit` as `Computed` ([#1423](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1423))
- datasource/gitlab_projects: Support validation for all valid `order_by` values ([#1429](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1429))

## 15.7.1 (2022-12-23)

This release was tested against GitLab 15.5, 15.6 and 15.7 for both CE and EE.

BUG FIXES:

- provider: default to `gitlab.com` as base URL if it is not provided ([#1400](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1400))

## 15.7.0 (2022-12-22)

- This release was tested against GitLab 15.5, 15.6 and 15.7 for both CE and EE.
- This is the first release aligned with GitLab. See [#1331](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/issues/1331).

BREAKING CHANGES:

See [Terraform GitLab Provider Version 15.7 Upgrade Guide](https://registry.terraform.io/providers/gitlabhq/gitlab/latest/docs/guides/version-15.7-upgrade) for details.

- Require at least Terraform 1.0 and Terraform Protocol Version 6 ([#1336](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1336))
- Provider `token` argument has changed to `sensitive` ([#1385](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1385))

FEATURES:

- resource/gitlab_project: Support forking a project ([#1377](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1377))
- resource/gitlab_project: Support avatars ([#1387](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1387))
- resource/gitlab_group: Support avatars ([#1387](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1387))
- **New Data Source**: `gitlab_metadata` ([#1355](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1355))

IMPROVEMENTS:

- resource/gitlab_project: Add `ci_separated_caches` attribute ([#1320](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1320))
- datasource/gitlab_project: Add `ci_separated_catches` attribute ([#1320](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1320))
- resource/gitlab_project: Rely on API defaults to only send minimal requests ([#1376](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1376))
- resource/gitlab_project: Add `restrict_user_defined_variables` attribute ([#1372](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1372))
- resource/gitlab_repository_file: Add `overwrite_on_create` attribute ([#1374](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1374))

BUG FIXES:

- resource/gitlab_user: Fix suppress logic for `skip_confirmation` attribute ([#1375](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1375))
- resource/gitlab_project: Fix disabling `container_expiration_policy` ([#1386](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1386))

## 3.20.0 (2022-11-25)

This release was tested against GitLab 15.4, 15.5 and 15.6 for both CE and EE.

FEATURES:

- **New Data Source:** `gitlab_groups` ([#1252](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1252))
- **New Data Source:** `gitlab_group_subgroups` ([#1280](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1280))
- **New Data Source:** `gitlab_user_sshkeys` ([#1296](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1296))

IMPROVEMENTS:

- resource/gitlab_project: deprecate `pipelines_enabled` and remove default ([#1357](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1357))
- resource/gitlab_project_access_token: Support `read_registry` and `write_registry` as valid scopes ([#1289](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1289))
- resource/gitlab_application_settings: Add `group_owners_can_manage_default_branch_protection` attribute ([#1334](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1334))
- resource/gitlab_user: Suppress undesired diff for `skip_confirmation` attribute ([#1339](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1339))
- resource/gitlab_group: Add `extra_shared_runners_minutes_limit` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))
- resource/gitlab_group: Add `membership_lock` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))
- resource/gitlab_group: Add `shared_runners_minutes_limit` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))
- datasource/gitlab_group: Add `extra_shared_runners_minutes_limit` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))
- datasource/gitlab_group: Add `membership_lock` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))
- datasource/gitlab_group: Add `shared_runners_minutes_limit` attribute ([#1232](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1232))

## 3.19.0 (2022-11-10)

This release was tested against GitLab 15.1, 15.2 and 15.3 for both CE and EE.

This is the first release from the [new project on GitLab](https://gitlab.com/gitlab-org/terraform-provider-gitlab).

FEATURES

- **New Resource:** `gitlab_service_emails_on_push` ([#1305](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1305))

IMPROVEMENTS:

- resource/gitlab_project: Add `suggestion_commit_message` attribute ([#1249](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1249))
- datasource/gitlab_project: Add `suggestion_commit_message` attribute ([#1249](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1249))
- datasource/gitlab_projects: Add `suggestion_commit_message` attribute ([#1249](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1249))
- resource/gitlab_branch_protection: Support `no one` in `unprotect_access_level` ([#1278](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1278))

BUG FIXES:

- resource/gitlab_repository_file: Fix check if file exists during read ([#1260](https://gitlab.com/gitlab-org/terraform-provider-gitlab/-/merge_requests/1260))

## 3.18.0 (2022-09-05)

This release was tested against GitLab 15.1, 15.2 and 15.3 for both CE and EE.

FEATURES:

- **New Resource:** `gitlab_group_saml_link` ([#1243](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1243))

## 3.17.0 (2022-08-24)

This release was tested against GitLab 15.1, 15.2 and 15.3 for both CE and EE.

FEATURES:

- **New Data Source:** `gitlab_project_hook` ([#1204](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1204))
- **New Data Source:** `gitlab_project_hooks` ([#1204](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1204))
- **New Data Source:** `gitlab_project_membership` ([#593](https://github.com/gitlabhq/terraform-provider-gitlab/pull/593))
- **New Data Source:** `gitlab_repository_tree` ([#1198](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1198))
- **New Data Source:** `gitlab_group_hook` ([#1221](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1221))
- **New Data Source:** `gitlab_group_hooks` ([#1221](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1221))
- **New Resource:** `gitlab_group_hook` ([#1221](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1221))
- **New Resource:** `gitlab_application_settings` (experimental) ([#1201](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1201))
- **New Resource:** `gitlab_project_issue_board` ([#1173](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1173))
- **New Resource:** `gitlab_user_gpgkey` ([#1181](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1181))

IMPROVEMENTS:

- resource/gitlab_service_jira: Support `jira_issue_transition_id` attribute for GitLab 15.2 ([#1188](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1188))
- resource/gitlab_project_protected_environment: Add docs that users and groups must be shared with the project ([#1210](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1210))
- resource/gitlab_branch_protection: Automatically take ownership of projects default branch without an import ([#1216](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1216))
- resource/gitlab_group_ldap_link: Clarify value for `ldap_provider` attribute ([#1220](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1220))
- resource/gitlab_group_membership: Support `skip_subresources_on_destroy` and `unassign_issuables_on_destroy` removal option attributes ([#1209](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1209))
- resource/gitlab_project: Reference doc for required `gitlab_group_project_template` resource when using `template_project_id` attribute ([#1223](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1223))

BUG FIXES:

- resource/gitlab_user_sshkey: Ignore leading and trailing whitespaces in key ([#1175](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1175))
- resource/gitlab_project: Fix setting `ci_forward_deployment_enabled` to `false` during creation ([#1218](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1218))

## 3.16.1 (2022-07-11)

This release was tested against GitLab 14.10, 15.0 and 15.1 for both CE and EE.

BUG FIXES:

- resource/gitlab_project: Fix admin token requirement to check default branch protection ([#1169](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1169))

## 3.16.0 (2022-07-07)

This release was tested against GitLab 14.10, 15.0 and 15.1 for both CE and EE.

FEATURES:

- **New Data Source:** `gitlab_current_user` ([#1118](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1118))
- **New Data Source:** `gitlab_release_link` ([#1131](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1131))
- **New Data Source:** `gitlab_release_links` ([#1131](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1131))
- **New Resource:** `gitlab_release_link` ([#1131](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1131))
- **New Resource:** `gitlab_cluster_agent_token` ([#1147](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1147))

IMPROVEMENTS:

- resource/gitlab_project_protected_environment: Add `required_approval_count` attribute ([#1097](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1097))
- resource/gitlab_project_access_token: Add `owner` as possible value to `access_level` ([#1145](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1145))
- resource/gitlab_project_membership: Add `owner` as possible value to `access_level` ([#1145](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1145))
- resource/gitlab_project_share_group: Add `owner` as possible value to `access_level` ([#1145](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1145))
- resource/gitlab_project: Add `ci_default_git_depth` attribute ([#1146](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1146))
- datasource/gitlab_project: Add `ci_default_git_depth` attribute ([#1146](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1146))
- datasource/gitlab_projects: Add `ci_default_git_depth` attribute ([#1146](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1146))

BUG FIXES:

- resource/gitlab_project: Fix project creation when default branch protection is disabled on instance-level ([#1128](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1128))
- resource/gitlab_project: Fix a case where a change to a project in terraform can never apply when certain fields are modified ([#1158](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1158))
- resource/gitlab_project: Fix passing `false` to API for explicitly set optional attributes ([#1152](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1152))
- resource/gitlab_group: Fix passing false to API for explicitly set optional attributes ([#1152](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1152))

## 3.15.1 (2022-06-08)

This release was tested against GitLab 14.9, 14.10 and 15.0 for both CE and EE.

BUG FIXES:

- resource/gitlab_service_microsoft_teams: Fix removal from state when integration is not found ([#1113](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1113))

## 3.15.0 (2022-05-29)

This release was tested against GitLab 14.9, 14.10 and 15.0 for both CE and EE.

FEATURES:

- **New Data Source:** `gitlab_cluster_agent` ([#1073](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1073))
- **New Data Source:** `gitlab_cluster_agents` ([#1073](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1073))
- **New Data Source:** `gitlab_project_milestone` ([#1044](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1044))
- **New Data Source:** `gitlab_project_milestones` ([#1044](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1044))
- **New Resource:** `gitlab_project_milestone` ([#1044](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1044))
- **New Resource:** `gitlab_runner` ([#1049](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1049))
- **New Resource:** `gitlab_cluster_agent` ([#1073](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1073))

IMPROVEMENTS:

- resource/gitlab_group: Allow value `3` for `default_branch_protection` attribute ([#1070](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1070))
- resource/gitlab_project_badge: Add `name` attribute ([#1052](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1052))
- resource/gitlab_group: Transfer a subgroup does not longer re-create the group ([#1078](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1078))
- resource/gitlab_topic: Add `name` attribute ([#1095](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1095))

BUG FIXES:

- resource/gitlab_project_issue: Remove `Optional` from `Computed`-only attributes ([#1081](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1081))
- datasource/gitlab_project_issues: Fix type of `not_milestone` attribute from list of strings to a single string ([#1095](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1095))

## 3.14.0 (2022-05-02)

FEATURES:

- **New Data Source:** `gitlab_project_variable` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Data Source:** `gitlab_project_variables` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Data Source:** `gitlab_group_variable` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Data Source:** `gitlab_group_variables` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Data Source:** `gitlab_instance_variable` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Data Source:** `gitlab_instance_variables` ([#990](https://github.com/gitlabhq/terraform-provider-gitlab/pull/990))
- **New Resource:** `gitlab_group_project_file_template` ([#971](https://github.com/gitlabhq/terraform-provider-gitlab/pull/971))
- **New Resource:** `gitlab_service_external_wiki` ([#1003](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1003))
- **New Resource:** `gitlab_project_runner_enablement` ([#1016](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1016))
- **New Resource:** `gitlab_personal_access_token` ([#1007](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1007))

IMPROVEMENTS:

- resource/gitlab_deploy_key: Fully support `can_push` attribute ([#1009](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1009))
- resource/gitlab_deploy_key_enable: Fully support `can_push` attribute ([#1009](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1009))
- resource/gitlab_group_access_token: Support `owner` as access level ([#999](https://github.com/gitlabhq/terraform-provider-gitlab/pull/999))
- resource/gitlab_pipeline_trigger: Mark `token` attribute as sensitive ([#1034](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1034))
- resource/gitlab_project: Deprecate `build_coverage_regex` ([#1036](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1036))
- resource/gitlab_project_access_token: Add `access_level` attribute ([#997](https://github.com/gitlabhq/terraform-provider-gitlab/pull/997))
- resource/gitlab_project_protected_environment: Support multiple `deploy_access_levels` ([#1004](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1004))
- resource/gitlab_project_mirror: Support deletion on destroy ([#988](https://github.com/gitlabhq/terraform-provider-gitlab/pull/988))
- resource/gitlab_repository_file: Add `execute_filemode` attribute ([#1038](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1038))
- resource/gitlab_user: Add `namespace_id` attribute ([#987](https://github.com/gitlabhq/terraform-provider-gitlab/pull/987))
- datasource/gitlab_user: Add `namespace_id` attribute ([#987](https://github.com/gitlabhq/terraform-provider-gitlab/pull/987))
- datasource/gitlab_users: Add `namespace_id` attribute ([#987](https://github.com/gitlabhq/terraform-provider-gitlab/pull/987))

BUG FIXES:

- resource/gitlab_service_slack: Fix a resource ID bug that causes `gitlab_service_slack` resources that were created before provider version 3.9.0 to return an error. ([#1013](https://github.com/gitlabhq/terraform-provider-gitlab/pull/1013))

## 3.13.0 (2022-03-30)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

FEATURES:

- **New Data Source:** `gitlab_instance_deploy_keys` ([#870](https://github.com/gitlabhq/terraform-provider-gitlab/pull/870))
- **New Data Source:** `gitlab_project_tags` ([#963](https://github.com/gitlabhq/terraform-provider-gitlab/pull/963))
- **New Data Source:** `gitlab_repository_file` ([#939](https://github.com/gitlabhq/terraform-provider-gitlab/pull/939))
- **New Resource**: `gitlab_project_environment` ([#938](https://github.com/gitlabhq/terraform-provider-gitlab/pull/938))
- **New Resource**: `gitlab_project_protected_environment` ([#938](https://github.com/gitlabhq/terraform-provider-gitlab/pull/938))
- **New Resource**: `gitlab_system_hook` ([#929](https://github.com/gitlabhq/terraform-provider-gitlab/pull/929))
- resource/gitlab_topic: Support deletion ([#967](https://github.com/gitlabhq/terraform-provider-gitlab/pull/967))
- resource/gitlab_topic: Support avatar images ([#968](https://github.com/gitlabhq/terraform-provider-gitlab/pull/968))
- resource/gitlab_repository_file: Support using plain text `content` to beautify plans ([#972](https://github.com/gitlabhq/terraform-provider-gitlab/pull/972))

IMPROVEMENTS:

- resource/gitlab_branch_protection: Make `push_access_level` and `merge_access_level` optional ([#934](https://github.com/gitlabhq/terraform-provider-gitlab/pull/934))
- resource/gitlab_branch_protection: Add `unprotect_access_level` and `allowed_to_unprotect` attributes ([#934](https://github.com/gitlabhq/terraform-provider-gitlab/pull/934))
- resource/gitlab_pipeline_schedule: Use single GET API to read resource details to increase performance ([#955](https://github.com/gitlabhq/terraform-provider-gitlab/pull/955))
- resource/gitlab_project_variable: Use single GET API to read resource details to increase performance ([#953](https://github.com/gitlabhq/terraform-provider-gitlab/pull/953))
- resource/gitlab_deploy_token: Support `terraform import` ([#960](https://github.com/gitlabhq/terraform-provider-gitlab/pull/960))
- resource/gitlab_project_access_token: Support `terraform import` ([#960](https://github.com/gitlabhq/terraform-provider-gitlab/pull/960))
- resource/gitlab_project_hook: Support `terraform import` ([#960](https://github.com/gitlabhq/terraform-provider-gitlab/pull/960))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `analytics_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `auto_cancel_pending_pipelines` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `auto_devops_deploy_strategy` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `auto_devops_enabled` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `autoclose_referenced_issues` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `build_git_strategy` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `builds_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `container_expiration_policy` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `container_registry_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `emails_disabled` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `external_authorization_classification_label` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `forking_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `issues_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `merge_commit_template` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `merge_requests_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `operations_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `public_builds` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `repository_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `repository_storage` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `requirements_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `resolve_outdated_diff_discussions` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `security_and_compliance_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `squash_commit_template` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `topics` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_project, datasource/gitlab_project, datasource/gitlab_projects: Add `wiki_access_level` attribute ([#917](https://github.com/gitlabhq/terraform-provider-gitlab/pull/917))
- resource/gitlab_managed_license: Add support for "allowed" and "denied" to align with upcoming GitLab 15.0 deprecations ([#952](https://github.com/gitlabhq/terraform-provider-gitlab/pull/952))

BUG FIXES:

- resource/gitlab_deploy_token: Implement pagination when reading tokens to find all existing tokens ([#941](https://github.com/gitlabhq/terraform-provider-gitlab/pull/941))
- resource/gitlab_project_approval_rule: Implement pagination when reading approval rules to find all existing ones ([#950](https://github.com/gitlabhq/terraform-provider-gitlab/pull/950))
- resource/gitlab_repository_file: Implement locking and retry within the provider to mitigate parallelism limits ([#964](https://github.com/gitlabhq/terraform-provider-gitlab/pull/964))

## 3.12.0 (2022-03-08)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

FEATURES:

- **New Resource:** `gitlab_project_tag` ([#910](https://github.com/gitlabhq/terraform-provider-gitlab/pull/910))
- **New Resource:** `gitlab_project_issue` ([#891](https://github.com/gitlabhq/terraform-provider-gitlab/pull/891))
- **New Data Source:** `gitlab_project_tag` ([#910](https://github.com/gitlabhq/terraform-provider-gitlab/pull/910))
- **New Data Source:** `gitlab_project_issue` ([#891](https://github.com/gitlabhq/terraform-provider-gitlab/pull/891))
- **New Data Source:** `gitlab_project_issues` ([#891](https://github.com/gitlabhq/terraform-provider-gitlab/pull/891))

IMPROVEMENTS:

- datasource/gitlab_group: Add `prevent_forking_outside_group` attribute ([#914](https://github.com/gitlabhq/terraform-provider-gitlab/pull/914))
- resource/gitlab_group: Add `prevent_forking_outside_group` attribute ([#914](https://github.com/gitlabhq/terraform-provider-gitlab/pull/914))
- resource/gitlab_project_approval_rule: Add `rule_type` attribute ([#916](https://github.com/gitlabhq/terraform-provider-gitlab/pull/916))

BUG FIXES:

- resource/gitlab_project: Fix deletion drift handling ([#924](https://github.com/gitlabhq/terraform-provider-gitlab/pull/924))
- resource/gitlab_project_badge: Fix deletion drift handling ([#924](https://github.com/gitlabhq/terraform-provider-gitlab/pull/924))
- resource/gitlab_group_badge: Fix deletion drift handling ([#924](https://github.com/gitlabhq/terraform-provider-gitlab/pull/924))

## 3.11.1 (2022-03-02)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

BUG FIXES:

- resource/gitlab_project: fix waiting for default branch protection during creation ([#908](https://github.com/gitlabhq/terraform-provider-gitlab/pull/908))

## 3.11.0 (2022-03-01)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

FEATURES:

- **New Resource:** `resource/gitlab_user_sshkey` ([#903](https://github.com/gitlabhq/terraform-provider-gitlab/pull/903))

IMPROVEMENTS:

- resource/gitlab_instance_variable: better error message for invalid masked variable values ([#895](https://github.com/gitlabhq/terraform-provider-gitlab/pull/895))
- resource/gitlab_group_variable: better error message for invalid masked variable values ([#895](https://github.com/gitlabhq/terraform-provider-gitlab/pull/895))
- resource/gitlab_project: Add `merge_pipelines_enabled` and `merge_trains_enabled` attributes ([#900](https://github.com/gitlabhq/terraform-provider-gitlab/pull/900))
- resource/gitlab_project_level_mr_approvals: Add `required_password_to_approve` attribute ([#808](https://github.com/gitlabhq/terraform-provider-gitlab/pull/808))
- resource/gitlab_user: Add support for `deactivated` user state ([#899](https://github.com/gitlabhq/terraform-provider-gitlab/pull/899))

BUG FIXES:

- resource/gitlab_branch_protection: fix issue claiming that no valid access level([#892](https://github.com/gitlabhq/terraform-provider-gitlab/pull/892))

## 3.10.1 (2022-02-24)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

BUG FIXES:

- resource/gitlab_branch_protection: Fix issue which claimed that `no valid access level` can be found ([#892](https://github.com/gitlabhq/terraform-provider-gitlab/pull/892))

## 3.10.0 (2022-02-23)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

FEATURES:

- **New Resource:** `gitlab_group_access_token` ([#860](https://github.com/gitlabhq/terraform-provider-gitlab/pull/860))
- **New Resource:** `gitlab_topic` ([#871](https://github.com/gitlabhq/terraform-provider-gitlab/pull/871))

IMPROVEMENTS:

- datasource/gitlab_project: Add `printing_merge_request_link_enabled` attribute ([#783](https://github.com/gitlabhq/terraform-provider-gitlab/pull/783))
- datasource/gitlab_project: Add `ci_forward_deployment_enabled` attribute ([#732](https://github.com/gitlabhq/terraform-provider-gitlab/pull/732))
- datasource/gitlab_projects: Add `ci_forward_deployment_enabled` attribute ([#732](https://github.com/gitlabhq/terraform-provider-gitlab/pull/732))
- datasource/gitlab_group_membership: Support pagination ([#858](https://github.com/gitlabhq/terraform-provider-gitlab/pull/858))
- resource/gitlab_group_ldap_link: Add import support ([#771](https://github.com/gitlabhq/terraform-provider-gitlab/pull/771))
- resource/gitlab_project: Add `ci_forward_deployment_enabled` attribute ([#732](https://github.com/gitlabhq/terraform-provider-gitlab/pull/732))
- resource/gitlab_project: Add `printing_merge_request_link_enabled` attribute ([#783](https://github.com/gitlabhq/terraform-provider-gitlab/pull/783))
- resource/gitlab_project_hook: Add `releases_events` attribute ([#773](https://github.com/gitlabhq/terraform-provider-gitlab/pull/773))
- resource/gitlab_branch_protection: Add `allow_force_push` attribute ([#877](https://github.com/gitlabhq/terraform-provider-gitlab/pull/877))
- resource/gitlab_service_jira: Add `api_url` attribute ([#597](https://github.com/gitlabhq/terraform-provider-gitlab/pull/597))
- resource/gitlab_user: Add `state` attribute to allow blocking users ([#762](https://github.com/gitlabhq/terraform-provider-gitlab/pull/762))

BUG FIXES:

- datasource/gitlab_projects: Allow to get archived and unarchived repositories ([#855](https://github.com/gitlabhq/terraform-provider-gitlab/pull/855))
- resource/gitlab_group: Support setting `default_branch_protection` to `0` ([#856](https://github.com/gitlabhq/terraform-provider-gitlab/pull/856))
- resource/gitlab_group_ldap_link: Fix panic when setting group access level ([#873](https://github.com/gitlabhq/terraform-provider-gitlab/pull/873))
- resource/gitlab_project: Correctly handle push rules add and edit ([#838](https://github.com/gitlabhq/terraform-provider-gitlab/pull/838))
- resource/gitlab_project: Support creating project in group without default branch protection ([#856](https://github.com/gitlabhq/terraform-provider-gitlab/pull/856))
- resource/gitlab_project: Fix backwards-compatibility with 14.1 regarding the `squash_option` ([#867](https://github.com/gitlabhq/terraform-provider-gitlab/pull/867))
- resource/gitlab_project: Re-compute `path_with_namespace`, `ssh_url_to_repo`, `http_url_to_repo` and `web_url` attributes if `path` changes ([#875](https://github.com/gitlabhq/terraform-provider-gitlab/pull/875))

## 3.9.1 (2022-02-06)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

BUGFIXES:

- Fix crash in `gitlab_user` data source

## 3.9.0 (2022-02-04)

NOTES:

- resource/gitlab_service_slack: There was a breaking change to this resource in version 3.9.0 which was fixed in 3.14.0. Thus, if you have `gitlab_service_slack` resources that were created before 3.9.0, you should skip this version.

FEATURES:

- 0519c96 - Add `gitlab_repository_file` resource (#724)
- a915ccb - Add the `gitlab_project_access_token` resource (#588)
- 8564a07 - Add `gitlab_group_custom_attribute`, `gitlab_project_custom_attribute`, and `gitlab_user_custom_attribute` (#731)
- 39b0b6d - Add support for `gitlab_group_badge` resource (#673)
- dd0da2b - Implement configuration of the integration "Microsoft Teams" (#308) (#784)
- febe745 - Add `gitlab_project_protected_branch` and `gitlab_project_protected_branches` data sources (#551)
- b4d4f8d - Support `early_auth_check` flag in provider config (#787)
- 1455db0 - Add `gitlab_managed_license` resource (#700)
- be59cd1 - New `archive_on_destroy` attribute for `gitlab_project` (#816)

ENHANCEMENTS:

- 9863a61 - Add missing attributes to documentation (#802)
- 84d729e - Improve documentation around authentication with GitLab (#803)
- a9edc4a - Add environment scope to `gitlab_group_variable` (#717)
- c5a2f34 - Support `issues_template` and `merge_requests_template` attributes in project resource (#796)
- 65b8e9c - Add minimal access level permissions (#747)
- 12ae791 - Added missing scopes to deploy token (#769)
- 1455db0 - Update allowed access levels in `gitlab_branch_protection`, `gitlab_group_ldap_link`, `gitlab_group_membership`, `gitlab_group_share_group`, `gitlab_project_membership`, `gitlab_project_share_group`, and `gitlab_tag_protection` (#819)
- 1455db0 - New attribute `group_access` in `gitlab_project_share_group` and `gitlab_group_ldap_link` (attribute `access_level` is now considered deprecated on those resources) (#819)

BUGFIXES:

- f3b0f31 - Handle manually removed resources (#318)
- fc84cc3 - Properly allow arguments `id`or `path_with_namespace` for project data source (#806)
- d7059cf - Do not recreate project mirror on every run (#632)
- e57bf1d - Fix branch_protection documentation (#780)
- 912b647 - Fix docs for `gitlab_project_level_mr_approvals` import (#766)
- e89922e - instance_cluster/group_cluster: Suppress whitespace diff for kubernetes_ca_cert (#728)
- edda219 - gitlab_user: Do not set skip_confirmation on read (#491)

## 3.8.0 (Nov 19, 2021)

ENHANCEMENTS:

- More examples and better documentation ([#654](https://github.com/gitlabhq/terraform-provider-gitlab/pull/654))
- Adding a first complete example code ([#670](https://github.com/gitlabhq/terraform-provider-gitlab/pull/670))
- Support running the provider using an OAuth token ([#686](https://github.com/gitlabhq/terraform-provider-gitlab/pull/686))
- Allow merge on skipped pipeline ([#705](https://github.com/gitlabhq/terraform-provider-gitlab/pull/705))
- support default_branch_protection on group ([#706](https://github.com/gitlabhq/terraform-provider-gitlab/pull/706))
- Bump gitlab.com/gitlab-org/api/client-go from 0.50.0 to 0.51.1 ([#718](https://github.com/gitlabhq/terraform-provider-gitlab/pull/718))
- Add squash option ([#719](https://github.com/gitlabhq/terraform-provider-gitlab/pull/719))

BUGFIXES:

- Suppress whitespace diffs for kubernetes_ca_cert attribute ([#665](https://github.com/gitlabhq/terraform-provider-gitlab/pull/665))
- Fix GitLab project membership example ([#684](https://github.com/gitlabhq/terraform-provider-gitlab/pull/684))
- Improve tests for initializing a project without README ([#730](https://github.com/gitlabhq/terraform-provider-gitlab/pull/730))

## 3.7.0 (July 20, 2021)

FEATURES:

- Add protected_branch_ids to gitlab_project_approval_rule ([#542](https://github.com/gitlabhq/terraform-provider-gitlab/pull/542))
- Add most Premium features for gitlab_branch_protection ([#556](https://github.com/gitlabhq/terraform-provider-gitlab/pull/556))
- Adds support for gitlab project badges ([#648](https://github.com/gitlabhq/terraform-provider-gitlab/pull/648))

ENHANCEMENTS:

- Add CustomCIPath to resource gitlab_project ([#662](https://github.com/gitlabhq/terraform-provider-gitlab/pull/662))
- add build coverage regex ([#627](https://github.com/gitlabhq/terraform-provider-gitlab/pull/627))
- Add tfproviderlint linter to CI ([#653](https://github.com/gitlabhq/terraform-provider-gitlab/pull/653))
- Documentation improvements ([#642](https://github.com/gitlabhq/terraform-provider-gitlab/pull/642))
- chore: add error handling to resourceGitlabProjectSetToState ([#649](https://github.com/gitlabhq/terraform-provider-gitlab/pull/649))
- Missing documentation for gitlab_instance_variable ([#602](https://github.com/gitlabhq/terraform-provider-gitlab/pull/602))
- Add support for importing Pipeline Schedules and Triggers ([#618](https://github.com/gitlabhq/terraform-provider-gitlab/pull/618))
- update reference to master branch to main ([#612](https://github.com/gitlabhq/terraform-provider-gitlab/pull/612))

BUGFIXES:

- Fix project approval rule tests ([#660](https://github.com/gitlabhq/terraform-provider-gitlab/pull/660))
- Fix markdown linter errors in branch_protection.md ([#656](https://github.com/gitlabhq/terraform-provider-gitlab/pull/656))
- fix: update default branch name to "main" ([#643](https://github.com/gitlabhq/terraform-provider-gitlab/pull/643))
- gitlab_project: Wait for template projects to be cloned/imported ([#621](https://github.com/gitlabhq/terraform-provider-gitlab/pull/621))

## 3.6.0 (Apr 6, 2021)

ENHANCEMENTS:

- Support the Note field in the User resource/data ([#572](https://github.com/gitlabhq/terraform-provider-gitlab/pull/572))
- add diffSuppressFunc for 'expires_at' attribute in 'deploy_token' ([#575](https://github.com/gitlabhq/terraform-provider-gitlab/pull/575))
- Update to Go 1.16 and terraform-plugin-sdk 1.16 ([#579](https://github.com/gitlabhq/terraform-provider-gitlab/pull/579))
- Bump gitlab.com/gitlab-org/api/client-go from 0.44.0 to 0.46.0 ([#586](https://github.com/gitlabhq/terraform-provider-gitlab/pull/586))

BUG FIXES:

- Fix "Run failed" on forks ([#565](https://github.com/gitlabhq/terraform-provider-gitlab/pull/565))

## 3.5.0 (Feb 18, 2021)

FEATURES:

- Add resource for project freeze periods ([#516](https://github.com/gitlabhq/terraform-provider-gitlab/pull/516))

ENHANCEMENTS:

- Update go version and go-gitlab version ([#523](https://github.com/gitlabhq/terraform-provider-gitlab/pull/523))
- Support additional attributes in `gitlab_project_hook` ([#525](https://github.com/gitlabhq/terraform-provider-gitlab/pull/525))
- Link badges in README to proper workflows ([#527](https://github.com/gitlabhq/terraform-provider-gitlab/pull/527))
- gitlab_project: Check each push rule individually ([#531](https://github.com/gitlabhq/terraform-provider-gitlab/pull/531))
- Allow `full_path` in addition to `id` in gitlab_project data source ([#532](https://github.com/gitlabhq/terraform-provider-gitlab/pull/532))
- Update test fixtures for better usability ([#535](https://github.com/gitlabhq/terraform-provider-gitlab/pull/535))
- Check for state change on user delete ([#539](https://github.com/gitlabhq/terraform-provider-gitlab/pull/539))
- Increase gitlab_project import timeout ([#536](https://github.com/gitlabhq/terraform-provider-gitlab/pull/536))
- Add optional mirror options ([#554](https://github.com/gitlabhq/terraform-provider-gitlab/pull/554))
- Remove vendor folder ([#546](https://github.com/gitlabhq/terraform-provider-gitlab/pull/546))
- Add dependabot config ([#558](https://github.com/gitlabhq/terraform-provider-gitlab/pull/558))
- Fix EE tests actually running against CE ([#564](https://github.com/gitlabhq/terraform-provider-gitlab/pull/564))
- Fix EE test mounting license as a directory ([#568](https://github.com/gitlabhq/terraform-provider-gitlab/pull/568))

BUG FIXES:

- fix deploy_token expiration ([#510](https://github.com/gitlabhq/terraform-provider-gitlab/pull/510))
- Fix group_share_group nil pointer reference ([#555](https://github.com/gitlabhq/terraform-provider-gitlab/pull/555))

## 3.4.0 (Jan 14, 2021)

FEATURES:

- Support sharing a group with another group ([#511](https://github.com/gitlabhq/terraform-provider-gitlab/pull/511))
- Support Project Mirroring ([#512](https://github.com/gitlabhq/terraform-provider-gitlab/pull/512))

## 3.3.0 (Nov 30, 2020)

FEATURES:

- Support instance level CI variables ([#389](https://github.com/gitlabhq/terraform-provider-gitlab/pull/389))

ENHANCEMENTS

- Add the pages_access_level parameter ([#472](https://github.com/gitlabhq/terraform-provider-gitlab/pull/472))
- Do not fail when project member does not exist ([#473](https://github.com/gitlabhq/terraform-provider-gitlab/pull/473))
- Make the runners_token on the project secret ([#474](https://github.com/gitlabhq/terraform-provider-gitlab/pull/474))
- Fix nil pointer dereference importing gitlab_user ([#490](https://github.com/gitlabhq/terraform-provider-gitlab/pull/490))
- Fix unit and acceptance tests not running ([#495](https://github.com/gitlabhq/terraform-provider-gitlab/pull/495))

## 3.2.0 (Nov 20, 2020)

FEATURES:

- Project Approval Rules ([#250](https://github.com/gitlabhq/terraform-provider-gitlab/pull/https://github.com/gitlabhq/terraform-provider-gitlab/pull/250))

ENHANCEMENTS

- Documentation for expires_at ([#482](https://github.com/gitlabhq/terraform-provider-gitlab/pull/482))
- Update set-env github action command ([484](https://github.com/gitlabhq/terraform-provider-gitlab/pull/484))

## 3.1.0 (Oct 16, 2020)

ENHANCEMENTS:

- Enable custom UserAgent ([#451](https://github.com/gitlabhq/terraform-provider-gitlab/pull/451))
- gitlab_project_mirror: Mark URL as sensitive ([#458](https://github.com/gitlabhq/terraform-provider-gitlab/pull/458))
- Remove old-style variable interpolation ([#456](https://github.com/gitlabhq/terraform-provider-gitlab/pull/456))

BUG FIXES:

- add pagination for ListPipelineSchedules ([#454](https://github.com/gitlabhq/terraform-provider-gitlab/pull/454))

## 3.0.0 (Sept 23, 2020)

BREAKING CHANGES:

- Resource `gitlab_project_push_rules` has been removed. You now instead specify project push rules using the `push_rules` attribute on the `gitlab_project` resource.
- The `shared_with_groups` attribute has been removed from the `gitlab_project` resource (but not the data source). You may use the `gitlab_project_share_group` resource instead.

NOTES:

- If you are using the `environment_scope` attribute of `gitlab_project_variable` to manage multiple variables with the same key, it is recommended to use GitLab 13.4+. See [this related GitLab issue](https://gitlab.com/gitlab-org/gitlab/-/issues/9912) for older versions.
- The ID format of the `gitlab_project_variable` resource changed. The upgrade should be automatic.
- The default value of the `gitlab_project_variable` resource's `environment_scope` attribute has changed from `0` to `*`.

FEATURES:

- **New Data Source:** `gitlab_group_membership` ([#264](https://github.com/gitlabhq/terraform-provider-gitlab/issues/264))
- **New Resource:** `gitlab_instance_cluster` ([#367](https://github.com/gitlabhq/terraform-provider-gitlab/issues/367))
- **New Resource:** `gitlab_project_level_mr_approvals` ([#356](https://github.com/gitlabhq/terraform-provider-gitlab/issues/356))
- **New Resource:** `gitlab_project_mirror` ([#358](https://github.com/gitlabhq/terraform-provider-gitlab/issues/358))
- **New Resource:** `gitlab_service_pipelines_email` ([#375](https://github.com/gitlabhq/terraform-provider-gitlab/issues/375))

ENHANCEMENTS:

- data-source/gitlab_project: New attributes `packages_enabled`, `path_with_namespace` and `push_rules` ([#405](https://github.com/gitlabhq/terraform-provider-gitlab/issues/405), [#403](https://github.com/gitlabhq/terraform-provider-gitlab/issues/403), [#422](https://github.com/gitlabhq/terraform-provider-gitlab/issues/422))
- resource/gitlab_branch_protection: New `code_owner_approval_required` attribute ([#380](https://github.com/gitlabhq/terraform-provider-gitlab/issues/380))
- resource/gitlab_project: New attributes `packages_enabled`, `path_with_namespace`, and `push_rules` ([#405](https://github.com/gitlabhq/terraform-provider-gitlab/issues/405), [#403](https://github.com/gitlabhq/terraform-provider-gitlab/issues/403), [#422](https://github.com/gitlabhq/terraform-provider-gitlab/issues/422))
- resource/gitlab_group: New attributes `share_with_group_lock`, `project_creation_level`, `auto_devops_enabled`, `emails_disabled`, `mentions_disabled`, `subgroup_creation_level`, `require_two_factor_authentication`, and `two_factor_grace_period` ([#362](https://github.com/gitlabhq/terraform-provider-gitlab/issues/362))
- resource/gitlab_group: Automatically detect removal ([#267](https://github.com/gitlabhq/terraform-provider-gitlab/issues/267))
- resource/gitlab_group_label: Can now be imported ([#339](https://github.com/gitlabhq/terraform-provider-gitlab/issues/339))
- resource/gitlab_project: New `import_url` attribute ([#381](https://github.com/gitlabhq/terraform-provider-gitlab/issues/381))
- resource/gitlab_project_push_rules: Can now be imported ([#360](https://github.com/gitlabhq/terraform-provider-gitlab/issues/360))
- resource/gitlab_project_variable: Better error message when a masked variable fails validation ([#371](https://github.com/gitlabhq/terraform-provider-gitlab/issues/371))
- resource/gitlab_project_variable: Automatically detect removal ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
- resource/gitlab_service_jira: Automatically detect removal ([#337](https://github.com/gitlabhq/terraform-provider-gitlab/issues/337))
- resource/gitlab_user: The `email` attribute can be changed without forcing recreation ([#261](https://github.com/gitlabhq/terraform-provider-gitlab/issues/261))
- resource/gitlab_user: Require either the `password` or `reset_password` attribute to be set ([#262](https://github.com/gitlabhq/terraform-provider-gitlab/issues/262))

BUG FIXES:

- resource/gitlab_pipeline_schedule: Fix a rare error during deletion ([#364](https://github.com/gitlabhq/terraform-provider-gitlab/issues/364))
- resource/gitlab_pipeline_schedule_variable: Fix a rare error during deletion ([#364](https://github.com/gitlabhq/terraform-provider-gitlab/issues/364))
- resource/gitlab_project: Fix the `default_branch` attribute changing to `null` after first apply ([#343](https://github.com/gitlabhq/terraform-provider-gitlab/issues/343))
- resource/gitlab_project_share_group: Fix the `access_level` attribute not updating ([#421](https://github.com/gitlabhq/terraform-provider-gitlab/issues/421))
- resource/gitlab_project_share_group: Fix the share not working if the project is also managed ([#421](https://github.com/gitlabhq/terraform-provider-gitlab/issues/421))
- resource/gitlab_project_variable: Fix inconsistent reads for variables with non-unique keys ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
- resource/gitlab_project_variable: Change the default `environment_scope` from `0` to `*` ([#409](https://github.com/gitlabhq/terraform-provider-gitlab/issues/409))
- resource/gitlab_service_jira: Fix a rare state inconsistency problem during creation ([#363](https://github.com/gitlabhq/terraform-provider-gitlab/issues/363))
- resource/gitlab_user: Fix some attributes saving incorrectly in state ([#261](https://github.com/gitlabhq/terraform-provider-gitlab/issues/261))

## 2.11.0 (July 24, 2020)

ENHANCEMENTS:

- Improvements to resource `gitlab_user` import
  ([#340](https://github.com/gitlabhq/terraform-provider-gitlab/issues/340))

## 2.10.0 (June 09, 2020)

FEATURES:

- **New Resource:** `gitlab_service_github`
  ([#311](https://github.com/gitlabhq/terraform-provider-gitlab/issues/311))

ENHANCEMENTS:

- add attribute `remove_source_branch_after_merge` to projects
  ([#289](https://github.com/gitlabhq/terraform-provider-gitlab/issues/289))

BUGFIXES:

- fix for flaky `gitlab_group` tests
  ([#320](https://github.com/gitlabhq/terraform-provider-gitlab/issues/320))
- Creating custom skip function for group_ldap_link tests.
  ([#328](https://github.com/gitlabhq/terraform-provider-gitlab/issues/328))

## 2.9.0 (June 01, 2020)

FEATURES:

- **New DataSource:** `gitlab_projects`
  ([#279](https://github.com/gitlabhq/terraform-provider-gitlab/issues/279))
- **New Resource:** `gitlab_deploy_token`
  ([#284](https://github.com/gitlabhq/terraform-provider-gitlab/issues/284))

ENHANCEMENTS:

- Add `management_project_id` for Group and Project Clusters
  ([#301](https://github.com/gitlabhq/terraform-provider-gitlab/issues/301))

## 2.8.0 (May 28, 2020)

FEATURES:

- **New Resource:** `gitlab_group_ldap_link`
  ([#296](https://github.com/gitlabhq/terraform-provider-gitlab/issues/296),
  [#316](https://github.com/gitlabhq/terraform-provider-gitlab/issues/316))

ENHANCEMENTS:

- Update resource gitlab_group_label to read labels from all pages
  ([#302](https://github.com/gitlabhq/terraform-provider-gitlab/issues/302))
- Provide a way to specify client cert and key
  ([#315](https://github.com/gitlabhq/terraform-provider-gitlab/issues/315))

BUGFIXES:

- Increase MaxIdleConnsPerHost in http.Transport
  ([#305](https://github.com/gitlabhq/terraform-provider-gitlab/issues/305))

## 2.7.0 (May 20, 2020)

- Implement `masked` parameters for `gitlab_group_variable`
  ([#271](https://github.com/gitlabhq/terraform-provider-gitlab/issues/271))

## 2.6.0 (April 08, 2020)

ENHANCEMENTS:

- Add jira flags
  ([#274](https://github.com/gitlabhq/terraform-provider-gitlab/issues/274))

## 2.5.1 (April 06, 2020)

BUGFIXES:

- Support for soft-delete of groups and projects in Gitlab Enterprise Edition
  ([#282](https://github.com/gitlabhq/terraform-provider-gitlab/issues/282),
  [#283](https://github.com/gitlabhq/terraform-provider-gitlab/issues/283),
  [#285](https://github.com/gitlabhq/terraform-provider-gitlab/issues/285),
  [#291](https://github.com/gitlabhq/terraform-provider-gitlab/issues/291))

ENHANCEMENTS:

- Switched from Travis CI to Github Actions
  ([#216](https://github.com/gitlabhq/terraform-provider-gitlab/issues/216))

## 2.5.0 (December 05, 2019)

ENHANCEMENTS:

- Implement `lfs_enabled`, `request_access_enabled`, and `pipelines_enabled` parameters for `gitlab_project`
  ([#225](https://github.com/gitlabhq/terraform-provider-gitlab/pull/225),
  [#226](https://github.com/gitlabhq/terraform-provider-gitlab/pull/226),
  [#227](https://github.com/gitlabhq/terraform-provider-gitlab/pull/227))

BUGFIXES:

- Fix label support when there is more than 20 labels on a project
  ([#229](https://github.com/gitlabhq/terraform-provider-gitlab/pull/229))
- Enable `environment_scope` for `gitlab_project_variable` lookup
  ([#228](https://github.com/gitlabhq/terraform-provider-gitlab/pull/229))
- Fix users data source when there is more than 20 users returned
  ([#230](https://github.com/gitlabhq/terraform-provider-gitlab/pull/230))

## 2.4.0 (November 28, 2019)

FEATURES:

- **New Resource:** `gitlab_group_label` ([#186](https://github.com/gitlabhq/terraform-provider-gitlab/pull/186))
- **New Resource:** `gitlab_group_cluster`
  ([#178](https://github.com/gitlabhq/terraform-provider-gitlab/pull/178))
- **New Resource:** `gitlab_pipeline_schedule_variable`
  ([#204](https://github.com/gitlabhq/terraform-provider-gitlab/pull/204))

ENHANCEMENTS:

- Add `runners_token` to gitlab groups ([#218](https://github.com/gitlabhq/terraform-provider-gitlab/pull/218))
- Add `reset_password` to `gitlab_user` ([#127](https://github.com/gitlabhq/terraform-provider-gitlab/pull/127))
- Update `access_level` available values ([#220](https://github.com/gitlabhq/terraform-provider-gitlab/pull/220))
- Make read callbacks graceful for `gitlab_project_share_group`, `gitlab_branch_protection` and
  `gitlab_label` resources ([#223](https://github.com/gitlabhq/terraform-provider-gitlab/pull/223))

BUGFIXES:

- Fix state not being updated for `gitlab_branch_protection`
  ([#166](https://github.com/gitlabhq/terraform-provider-gitlab/pull/166))
- Set ForceNew for `gitlab_pipeline_schedule` `project`
  ([#203](https://github.com/gitlabhq/terraform-provider-gitlab/pull/203))

## 2.3.0 (October 17, 2019)

_We would like to thank Gitlab, which has provided us a EE license. This project
is now tested against Gitlab CE and Gitlab EE._

FEATURES:

- **New Resource:** `gitlab_project_push_rules` ([#163](https://github.com/gitlabhq/terraform-provider-gitlab/pull/163))
- **New Resource:** `gitlab_deploy_key_enable` ([#176](https://github.com/gitlabhq/terraform-provider-gitlab/pull/176))
- **New Resource:** `gitlab_project_share_group` ([#167](https://github.com/gitlabhq/terraform-provider-gitlab/pull/167))

ENHANCEMENTS:

- Add `initialize_with_readme` to `gitlab_project` ([#179](https://github.com/gitlabhq/terraform-provider-gitlab/issues/179))
- Add support for more variable options ([#169](https://github.com/gitlabhq/terraform-provider-gitlab/issues/169))
- Documentation improvements ([#168](https://github.com/gitlabhq/terraform-provider-gitlab/issues/168), [#187](https://github.com/gitlabhq/terraform-provider-gitlab/issues/187), [#171](https://github.com/gitlabhq/terraform-provider-gitlab/issues/171))

BUGFIXES:

- Fix tag protection URL
  ([#156](https://github.com/gitlabhq/terraform-provider-gitlab/issues/156))
- Properly manage the default branch in a git repo
  ([#158](https://github.com/gitlabhq/terraform-provider-gitlab/issues/158))
- Resolve triggers pagination issue by calling `GetPipelineTrigger`
  ([#173](https://github.com/gitlabhq/terraform-provider-gitlab/issues/173))

## 2.2.0 (June 12, 2019)

FEATURES:

- **New Resource:** `gitlab_service_jira` ([#101](https://github.com/gitlabhq/terraform-provider-gitlab/pull/101))
- **New Resource:** `gitlab_pipeline_schedule` ([#116](https://github.com/gitlabhq/terraform-provider-gitlab/pull/116))

ENHANCEMENTS:

- Add `archived` argument to `gitlab_project` ([#148](https://github.com/gitlabhq/terraform-provider-gitlab/issues/148))
- Add `managed` argument to `gitlab_project_cluster` ([#137](https://github.com/gitlabhq/terraform-provider-gitlab/issues/137))

## 2.1.0 (May 29, 2019)

FEATURES:

- **New Datasource**: `gitlab_group` ([#129](https://github.com/gitlabhq/terraform-provider-gitlab/issues/129))

## 2.0.0 (May 23, 2019)

This is the first release to support Terraform 0.12.

BACKWARDS INCOMPATIBILITIES:

- **all**: Previous versions of this provider silently removed state from state when
  Gitlab returned an error 404. Now we error on this and you must reconciliate
  the state (e.g. `terraform state rm`). We have done this because we can not
  make the difference between permission denied and resources removed outside of
  terraform (gitlab returns 404 in both cases)
  ([#130](https://github.com/gitlabhq/terraform-provider-gitlab/pull/130))

FEATURES:

- **New Resource:** `gitlab_tag_protection` ([#125](https://github.com/gitlabhq/terraform-provider-gitlab/pull/125))

ENHANCEMENTS:

- Add `container_registry_enabled` argument to `gitlab_project` ([#115](https://github.com/gitlabhq/terraform-provider-gitlab/issues/115))
- Add `shared_runners_enabled` argument to `gitlab_project` ([#134](https://github.com/gitlabhq/terraform-provider-gitlab/issues/134) [#104](https://github.com/gitlabhq/terraform-provider-gitlab/issues/104))

## 1.3.0 (May 03, 2019)

FEATURES:

- **New Resource:** `gitlab_service_slack` ([#96](https://github.com/gitlabhq/terraform-provider-gitlab/issues/96))
- **New Resource:** `gitlab_branch_protection` ([#68](https://github.com/gitlabhq/terraform-provider-gitlab/issues/68))

ENHANCEMENTS:

- Support for request/response logging when >`DEBUG` severity is set ([#93](https://github.com/gitlabhq/terraform-provider-gitlab/issues/93))
- Datasource `gitlab_user` supports user_id, email lookup and return lots of new attributes ([#102](https://github.com/gitlabhq/terraform-provider-gitlab/issues/102))
- Resource `gitlab_deploy_key` can now be imported ([#197](https://github.com/gitlabhq/terraform-provider-gitlab/issues/97))
- Add `tags` attribute for `gitlab_project` ([#106](https://github.com/gitlabhq/terraform-provider-gitlab/pull/106))

BUGFIXES:

- Documentation fixes ([#108](https://github.com/gitlabhq/terraform-provider-gitlab/issues/108), [#113](https://github.com/gitlabhq/terraform-provider-gitlab/issues/113))

## 1.2.0 (February 19, 2019)

FEATURES:

- **New Datasource:** `gitlab_users` ([#79](https://github.com/gitlabhq/terraform-provider-gitlab/issues/79))
- **New Resource:** `gitlab_pipeline_trigger` ([#82](https://github.com/gitlabhq/terraform-provider-gitlab/issues/82))
- **New Resource:** `gitlab_project_cluster` ([#87](https://github.com/gitlabhq/terraform-provider-gitlab/issues/87))

ENHANCEMENTS:

- Supports "No one" and "maintainer" permissions ([#83](https://github.com/gitlabhq/terraform-provider-gitlab/issues/83))
- `gitlab_project.shared_with_groups` is now order-independent ([#86](https://github.com/gitlabhq/terraform-provider-gitlab/issues/86))
- add `merge_method`, `only_allow_merge_if_*`, `approvals_before_merge` parameters to `gitlab_project` ([#72](https://github.com/gitlabhq/terraform-provider-gitlab/issues/72), [#88](https://github.com/gitlabhq/terraform-provider-gitlab/issues/88))

## 1.1.0 (January 14, 2019)

FEATURES:

- **New Resource:** `gitlab_project_membership`
- **New Resource:** `gitlab_group_membership` ([#8](https://github.com/gitlabhq/terraform-provider-gitlab/issues/8))
- **New Resource:** `gitlab_project_variable` ([#47](https://github.com/gitlabhq/terraform-provider-gitlab/issues/47))
- **New Resource:** `gitlab_group_variable` ([#47](https://github.com/gitlabhq/terraform-provider-gitlab/issues/47))

BACKWARDS INCOMPATIBILITIES:

`gitlab_project_membership` is not compatible with a previous _unreleased_ version due to an id change resource will need to be reimported manually
e.g

```bash
terraform state rm gitlab_project_membership.foo
terraform import gitlab_project_membership.foo 12345:1337
```

## 1.0.0 (October 06, 2017)

BACKWARDS INCOMPATIBILITIES:

- This provider now uses the v4 api. It means that if you set up a custom API url, you need to update it to use the /api/v4 url. As a side effect, we no longer support Gitlab < 9.0. ([#20](https://github.com/gitlabhq/terraform-provider-gitlab/issues/20))
- We now support Parent ID for `gitlab_groups`. However, due to a limitation in
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

- **New Resource:** `gitlab_user` ([#23](https://github.com/gitlabhq/terraform-provider-gitlab/issues/23))
- **New Resource:** `gitlab_label` ([#22](https://github.com/gitlabhq/terraform-provider-gitlab/issues/22))

IMPROVEMENTS:

- Add `cacert_file` and `insecure` options to the provider. ([#5](https://github.com/gitlabhq/terraform-provider-gitlab/issues/5))
- Fix race conditions with `gitlab_project` deletion. ([#19](https://github.com/gitlabhq/terraform-provider-gitlab/issues/19))
- Add `parent_id` argument to `gitlab_group`. ([#28](https://github.com/gitlabhq/terraform-provider-gitlab/issues/28))
- Add support for `gitlab_project` import. ([#30](https://github.com/gitlabhq/terraform-provider-gitlab/issues/30))
- Add support for `gitlab_groups` import. ([#31](https://github.com/gitlabhq/terraform-provider-gitlab/issues/31))
- Add `path` argument for `gitlab_project`. ([#21](https://github.com/gitlabhq/terraform-provider-gitlab/issues/21))
- Fix indempotency issue with `gitlab_deploy_key` and white spaces. ([#34](https://github.com/gitlabhq/terraform-provider-gitlab/issues/34))

## 0.1.0 (June 20, 2017)

NOTES:

- Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
