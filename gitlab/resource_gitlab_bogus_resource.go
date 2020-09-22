package gitlab

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

func resourceGitLabBogusResource() *schema.Resource {
	return &schema.Resource{
		Create: func(_ *schema.ResourceData, _ interface{}) error { return nil },
		Read:   func(_ *schema.ResourceData, _ interface{}) error { return nil },
		Delete: func(_ *schema.ResourceData, _ interface{}) error { return nil },
		Schema: map[string]*schema.Schema{
			"foo": {Type: schema.TypeInt, Required: true, ForceNew: true},
		},
	}
}
