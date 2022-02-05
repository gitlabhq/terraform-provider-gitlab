package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

var _ = registerResource("gitlab_group_custom_attribute", func() *schema.Resource {
	return CreateCustomAttributeResource(
		"group",
		func(client *gitlab.Client) CustomAttributeGetter {
			return client.CustomAttribute.GetCustomGroupAttribute
		},
		func(client *gitlab.Client) CustomAttributeSetter {
			return client.CustomAttribute.SetCustomGroupAttribute
		},
		func(client *gitlab.Client) CustomAttributeDeleter {
			return client.CustomAttribute.DeleteCustomGroupAttribute
		},
		"This resource allows you to set custom attributes for a group.",
	)
})
