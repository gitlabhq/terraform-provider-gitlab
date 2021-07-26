# gitlab\_project\_hook

This resource allows you to create and manage hooks for your GitLab projects.
For further information on hooks, consult the [gitlab
documentation](https://docs.gitlab.com/ce/user/project/integrations/webhooks.html).

## Example Usage

```hcl
resource "gitlab_project_hook" "example" {
  project               = "example/hooked"
  url                   = "https://example.com/hook/example"
  merge_requests_events = true
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) The name or id of the project to add the hook to.

* `url` - (Required) The url of the hook to invoke.

* `token` - (Optional) A token to present when invoking the hook.

* `enable_ssl_verification` - (Optional) Enable ssl verification when invoking the hook.

* `push_events` - (Optional) Invoke the hook for push events.

* `push_events_branch_filter` - (Optional) Invoke the hook for push events on matching branches only.

* `issues_events` - (Optional) Invoke the hook for issues events.

* `confidential_issues_events` - (Optional) Invoke the hook for confidential issues events.

* `merge_requests_events` - (Optional) Invoke the hook for merge requests.

* `tag_push_events` - (Optional) Invoke the hook for tag push events.

* `note_events` - (Optional) Invoke the hook for notes events.

* `confidential_note_events` - (Optional) Invoke the hook for confidential notes events.

* `job_events` - (Optional) Invoke the hook for job events.

* `pipeline_events` - (Optional) Invoke the hook for pipeline events.

* `wiki_page_events` - (Optional) Invoke the hook for wiki page events.
  
* `deployment_events` - (Optional) Invoke the hook for deployment events.

## Attributes Reference

The resource exports the following attributes:

* `id` - The unique id assigned to the hook by the GitLab server.
