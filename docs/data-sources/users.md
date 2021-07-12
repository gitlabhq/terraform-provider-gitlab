# gitlab\_users

Provide details about a list of users in the gitlab provider. The results include id, username, email, name and more about the requested users. Users can also be sorted and filtered using several options.

**NOTE**: Some available options require administrator privileges. Please visit [Gitlab API documentation][users_for_admins] for more information.

## Example Usage

```hcl
data "gitlab_users" "example" {
  sort = "desc"
  order_by = "name"
  created_before = "2019-01-01"
}

data "gitlab_users" "example-two" {
  search = "username"
}
```

## Argument Reference

The following arguments are supported:

* `search` - (Optional) Search users by username, name or email.

* `active` - (Optional) Filter users that are active.

* `blocked` - (Optional) Filter users that are blocked.

* `order_by` - (Optional) Order the users' list by `id`, `name`, `username`, `created_at` or `updated_at`. (Requires administrator privileges)

* `sort` - (Optional) Sort users' list in asc or desc order. (Requires administrator privileges)

* `extern_uid` - (Optional) Lookup users by external UID. (Requires administrator privileges)

* `extern_provider` - (Optional) Lookup users by external provider. (Requires administrator privileges)

* `created_before` - (Optional) Search for users created before a specific date. (Requires administrator privileges)

* `created_after` - (Optional)  Search for users created after a specific date. (Requires administrator privileges)

## Attributes Reference

The following attributes are exported:

* `users` - The list of users.
  * `id` - The unique id assigned to the user by the gitlab server.
  * `username` - The username of the user.
  * `email` - The e-mail address of the user.
  * `name` - The name of the user.
  * `is_admin` - Whether the user is an admin.
  * `can_create_group` - Whether the user can create groups.
  * `can_create_project` - Whether the user can create projects.
  * `projects_limit` - Number of projects the user can create.
  * `created_at` - Date the user was created at.
  * `state` - Whether the user is active or blocked.
  * `external` - Whether the user is external.
  * `extern_uid` - The external UID of the user.
  * `provider` - The UID provider of the user.
  * `organization` - The organization of the user.
  * `two_factor_enabled` - Whether user's two-factor auth is enabled.
  * `note` - The note associated to the user.
  * `avatar_url` - The avatar URL of the user.
  * `bio` - The bio of the user.
  * `location` - The location of the user.
  * `skype` - Skype username of the user.
  * `linkedin` - LinkedIn profile of the user.
  * `twitter` - Twitter username of the user.
  * `website_url` - User's website URL.
  * `theme_id` - User's theme ID.
  * `color_scheme_id` - User's color scheme ID.
  * `last_sign_in_at` - Last user's sign-in date.
  * `current_sign_in_at` - Current user's sign-in date.

[users_for_admins]: https://docs.gitlab.com/ce/api/users.html#for-admins
