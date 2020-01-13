---
layout: "gitlab"
page_title: "GitLab: gitlab_merge_request_approvals"
sidebar_current: "docs-gitlab-resource-label"
description: |-
  Creates and merge request approvals for GitLab projects
---

# gitlab\_project\_approvals\_configuration

This resource allows you to create and merge request approvals for your GitLab projects.
For further information on merge request approvals, consult the [gitlab
documentation](https://docs.gitlab.com/ee/api/merge_request_approvals.html).


## Example Usage

```hcl
resource "gitlab_project" "example" {
  name              = "example"
  description       = "My awesome codebase"
  visibility_level  = "public"
}

resource "gitlab_project_approvals_configuration" "approvals" {
  project_id                                     = gitlab_project.example.id
  reset_approvals_on_push                        = true
  disable_overriding_approvers_per_merge_request = true
  merge_requests_author_approval                 = false
  merge_requests_disable_committers_approval     = false
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The id of the project to add the label to.

* `reset_approvals_on_push` - (Optional) If set to true number of approvals is set to 0 after each push to the merge request.

* `disable_overriding_approvers_per_merge_request` - (Optional) If set to true it is not allowed to change project settings for merge request approvers in single merge request.

* `merge_requests_author_approval` - (Optional) If set to true author of merge request can approve it.

* `merge_requests_disable_committers_approval` - (Optional) If set to true commiters of merge request can approve it.

## Importing approvals configuration

You can import a group state using `terraform import <resource> <project_id>`. For example:

    terraform import gitlab_project_approvals_configuration.approvals 1117028
