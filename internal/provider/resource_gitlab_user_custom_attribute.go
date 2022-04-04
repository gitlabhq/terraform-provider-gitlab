package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_user_custom_attribute", func() *schema.Resource {
	return CreateCustomAttributeResource(
		"user",
		func(client *gitlab.Client) CustomAttributeGetter {
			return client.CustomAttribute.GetCustomUserAttribute
		},
		func(client *gitlab.Client) CustomAttributeSetter {
			return client.CustomAttribute.SetCustomUserAttribute
		},
		func(client *gitlab.Client) CustomAttributeDeleter {
			return client.CustomAttribute.DeleteCustomUserAttribute
		},
		`The `+"`gitlab_user_custom_attribute`"+` resource allows to manage custom attributes for a user.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/custom_attributes.html)`,
	)
})
