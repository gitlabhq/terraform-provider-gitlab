# gitlab\_service\_github

**NOTE**: requires either EE (self-hosted) or Silver and above (GitLab.com).

This resource manages a [GitHub integration](https://docs.gitlab.com/ee/user/project/integrations/github.html) that updates pipeline statuses on a GitHub repo's pull requests.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name = "awesome_project"
  description = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_github" "github" {
  project        = gitlab_project.awesome_project.id
  token          = "REDACTED"
  repository_url = "https://github.com/gitlabhq/terraform-provider-gitlab"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required) ID of the project you want to activate integration on.

* `repository_url` - (Required) The URL of the GitHub repo to integrate with, e,g, https://github.com/gitlabhq/terraform-provider-gitlab.

* `token` - (Required) A GitHub personal access token with at least `repo:status` scope.

* `static_context` - (Optional) Append instance name instead of branch to the status. Must enable to set a GitLab status check as _required_ in GitHub. See [Static / dynamic status check names] to learn more.

## Import

 You can import a service_github state using `terraform import <resource> <project_id>`:

```
$ terraform import gitlab_service_github.github 1
```
