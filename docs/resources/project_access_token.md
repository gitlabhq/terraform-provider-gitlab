# gitlab\_project\_access\_token

This resource allows you to create and manage Project Access Token for your GitLab projects.


## Example Usage

```hcl
resource "gitlab_deploy_token" "example" {
  project    = "example/deploying"
  name       = "Example project access token"
  expires_at = "2020-03-14"
  
  scopes = [ "api" ]
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Required, string) The id of the project to add the deploy token to.
 
* `name` - (Required, string) A name to describe the deploy token with.

* `expires_at` - (Optional, string) Time the token will expire it, YYYY-MM-DD format. Will not expire per default.

* `scopes` - (Required, set of strings) Valid values: `api`, `read_api`, `read_repository`, `write_repository`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `token` - The secret token. This is only populated when creating a new deploy token.

* `active` - True if the token is active.

* `created_at` - Time the token has been created, RFC3339 format.

* `revoked` - True if the token is revoked.

* `user_id` - The user_id associated to the token.
