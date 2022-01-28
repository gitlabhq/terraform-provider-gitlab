package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	gitlab "github.com/xanzy/go-gitlab"
)

func resourceGitlabProjectCustomAttribute() *schema.Resource {
	return CreateCustomAttributeResource(
		"project",
		func(client *gitlab.Client) CustomAttributeGetter {
			return client.CustomAttribute.GetCustomProjectAttribute
		},
		func(client *gitlab.Client) CustomAttributeSetter {
			return client.CustomAttribute.SetCustomProjectAttribute
		},
		func(client *gitlab.Client) CustomAttributeDeleter {
			return client.CustomAttribute.DeleteCustomProjectAttribute
		},
		"This resource allows you to set custom attributes for a project.",
	)
}
