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
		"This resource allows you to set custom attributes for a user.",
	)
})
