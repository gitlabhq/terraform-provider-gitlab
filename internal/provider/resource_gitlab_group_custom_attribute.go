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
		`The `+"`gitlab_group_custom_attribute`"+` resource allows to manage custom attributes for a group.

**Upstream API**: [GitLab REST API docs](https://docs.gitlab.com/ee/api/custom_attributes.html)`,
	)
})
