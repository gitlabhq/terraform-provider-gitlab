---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "gitlab_release_links Data Source - terraform-provider-gitlab"
subcategory: ""
description: |-
  The gitlab_release_links data source allows get details of release links.
  Upstream API: GitLab REST API docs https://docs.gitlab.com/api/releases/links/
---

# gitlab_release_links (Data Source)

The `gitlab_release_links` data source allows get details of release links.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/api/releases/links/)

## Example Usage

```terraform
# By project ID
data "gitlab_release_links" "example" {
  project  = "12345"
  tag_name = "v1.0.1"
}

# By project full path
data "gitlab_release_links" "example" {
  project  = "foo/bar"
  tag_name = "v1.0.1"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `project` (String) The ID or full path to the project.
- `tag_name` (String) The tag associated with the Release.

### Read-Only

- `id` (String) The ID of this resource.
- `release_links` (List of Object) List of release links (see [below for nested schema](#nestedatt--release_links))

<a id="nestedatt--release_links"></a>
### Nested Schema for `release_links`

Read-Only:

- `direct_asset_url` (String)
- `external` (Boolean)
- `filepath` (String)
- `link_id` (Number)
- `link_type` (String)
- `name` (String)
- `project` (String)
- `tag_name` (String)
- `url` (String)
