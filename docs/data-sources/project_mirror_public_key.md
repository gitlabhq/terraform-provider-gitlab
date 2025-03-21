---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_project_mirror_public_key Data Source - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_project_mirror_public_key data source allows the public key of a project mirror to be retrieved by its mirror id and the project it belongs to.
  Note: Supported on GitLab 17.9 or higher.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/remote_mirrors/#get-a-single-projects-remote-mirror-public-key
---

# gitlab_project_mirror_public_key (Data Source)

The `gitlab_project_mirror_public_key` data source allows the public key of a project mirror to be retrieved by its mirror id and the project it belongs to.

**Note**: Supported on GitLab 17.9 or higher.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/remote_mirrors/#get-a-single-projects-remote-mirror-public-key)

## Example Usage

```terraform
data "gitlab_project_mirror_public_key" "example" {
  project_id = 30
  mirror_id  = 42
}

data "gitlab_project_mirror_public_key" "example" {
  project_id = "foo/bar/baz"
  mirror_id  = 123
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `mirror_id` (Number) The id of the remote mirror.
- `project_id` (String) The integer or path with namespace that uniquely identifies the project.

### Read-Only

- `id` (String) The ID of this Terraform resource. In the format of `<project_id>:<mirror_id>`.
- `public_key` (String) Public key of the remote mirror.
