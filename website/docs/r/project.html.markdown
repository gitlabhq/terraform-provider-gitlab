---
layout: "gitlab"
page_title: "GitLab: gitlab_project"
sidebar_current: "docs-gitlab-resource-project-x"
description: |-
  Creates and manages projects within GitLab groups or within your user
---

# gitlab\_project

This resource allows you to create and manage projects within your
GitLab group or within your user.


## Example Usage

```hcl
resource "gitlab_project" "example" {
  name        = "example"
  description = "My awesome codebase"

  visibility_level = "public"
}
```

### Template Examples

Project templates can pre-populate a new project with the necessary files to get you started quickly.

### Built-In Template Project

Gitlab has [built-in templates](https://docs.gitlab.com/ee/gitlab-basics/create-project.html#built-in-templates), sourced from the following groups:
* [project-templates](https://gitlab.com/gitlab-org/project-templates)
* [pages](https://gitlab.com/pages)


```hcl
resource "gitlab_project" "template" {
  name          = "template"
  template_name = "gatsby"
}
```

###  Custom Template Project

GitLab administrators and users can configure their own [Custom project templates](https://docs.gitlab.com/ee/gitlab-basics/create-project.html#custom-project-templates-premium) 

```hcl
resource "gitlab_project" "custom_template" {
  name = "my-new-project"
  template_name = "my-custom-template"
  use_custom_template= true
  group_with_project_templates_id = 1
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the project.

* `path` - (Optional) The path of the repository.

* `namespace_id` - (Optional) The namespace (group or user) of the project. Defaults to your user.
  See [`gitlab_group`](group.html) for an example.

* `description` - (Optional) A description of the project.

* `tags` - (Optional) Tags (topics) of the project.

* `default_branch` - (Optional) The default branch for the project.

* `request_access_enabled` - Allow users to request member access.

* `issues_enabled` - (Optional) Enable issue tracking for the project.

* `merge_requests_enabled` - (Optional) Enable merge requests for the project.

* `pipelines_enabled` - (Optional) Enable pipelines for the project.

* `approvals_before_merge` - (Optional) Number of merge request approvals required for merging. Default is 0.

* `wiki_enabled` - (Optional) Enable wiki for the project.

* `snippets_enabled` - (Optional) Enable snippets for the project.

* `container_registry_enabled` - (Optional) Enable container registry for the project.

* `lfs_enabled` - (Optional) Enable LFS for the project.

* `visibility_level` - (Optional) Set to `public` to create a public project.
  Valid values are `private`, `internal`, `public`.
  Repositories are created as private by default.

* `merge_method` - (Optional) Set to `ff` to create fast-forward merges
  Valid values are `merge`, `rebase_merge`, `ff`
  Repositories are created with `merge` by default

* `only_allow_merge_if_pipeline_succeeds` - (Optional) Set to true if you want allow merges only if a pipeline succeeds.

* `only_allow_merge_if_all_discussions_are_resolved` - (Optional) Set to true if you want allow merges only if all discussions are resolved.

* `shared_runners_enabled` - (Optional) Enable shared runners for this project.

* `shared_with_groups` - (Optional) Enable sharing the project with a list of groups (maps).
  * `group_id` - (Required) Group id of the group you want to share the project with.
  * `group_access_level` - (Required) Group's sharing permissions. See [group members permission][group_members_permissions] for more info.
  Valid values are `guest`, `reporter`, `developer`, `master`.

* `archived` - (Optional) Whether the project is in read-only mode (archived). Repositories can be archived/unarchived by toggling this parameter.

* `initialize_with_readme` - (Optional) Create master branch with first commit containing a README.md file.

* `template_name` - (Optional)  Create project from existing template. When used without `use_custom_template`, name of a built-in project template. When used with `use_custom_template`, name of a custom project template

* `use_custom_template` - (Optional) Use either custom instance or group (with `group_with_project_templates_id`) project template

* `group_with_project_templates_id` - (Optional) For group-level custom templates, specifies ID of group from which all the custom project templates are sourced. Leave empty for instance-level templates. Requires use_custom_template to be true
## Attributes Reference

The following additional attributes are exported:

* `id` - Integer that uniquely identifies the project within the gitlab install.

* `ssh_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via SSH.

* `http_url_to_repo` - URL that can be provided to `git clone` to clone the
  repository via HTTP.

* `web_url` - URL that can be used to find the project in a browser.

* `runners_token` - Registration token to use during runner setup.

* `shared_with_groups` - List of the groups the project is shared with.
  * `group_name` - Group's name.

* `remove_source_branch_after_merge` - Enable `Delete source branch` option by default for all new merge requests.


## Importing projects

You can import a project state using `terraform import <resource> <id>`.  The
`id` can be whatever the [get single project api][get_single_project] takes for
its `:id` value, so for example:

    terraform import gitlab_project.example richardc/example

[get_single_project]: https://docs.gitlab.com/ee/api/projects.html#get-single-project
[group_members_permissions]: https://docs.gitlab.com/ce/user/permissions.html#group-members-permissions
