# gitlab\_service\_pipelines_email

This resource manages a [Pipelines email integration](https://docs.gitlab.com/ee/user/project/integrations/overview.html#integrations-listing) that emails the pipeline status to a list of recipients.

## Example Usage

```hcl
resource "gitlab_project" "awesome_project" {
  name             = "awesome_project"
  description      = "My awesome project."
  visibility_level = "public"
}

resource "gitlab_service_pipelines_email" "email" {
  project                      = gitlab_project.awesome_project.id
  recipients                   = ["gitlab@user.create"]
  notify_only_broken_pipelines = true
  branches_to_be_notified      = "all"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) ID of the project you want to activate integration on.

* `recipients` - (Required, set(string)) email addresses where notifications are sent.

* `notify_only_broken_pipelines` - (Optional, bool) Notify only broken pipelines. Default is true.

* `branches_to_be_notified` - (Optional, string) Branches to send notifications for. Valid options are `all`, `default`, `protected`, and `default_and_protected`. Default is `default`

## Importing Pipelines email service

 You can import a service_pipelines_email state using `terraform import <resource> <project_id>`:

```bash
$ terraform import gitlab_service_pipelines_email.email 1
```
