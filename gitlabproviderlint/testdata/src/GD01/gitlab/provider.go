package gitlab

import "github.com/hashicorp/terraform-plugin-sdk/helper/schema"

var _ = schema.Provider{
	DataSourcesMap: map[string]*schema.Resource{
		"gitlab_nice_datasource":    nil,
		"gitlab_naughty_datasource": nil, // want `Data source "gitlab_naughty_datasource" is missing a docs page named "naughty_datasource.md"`
	},
	ResourcesMap: map[string]*schema.Resource{
		"gitlab_nice_resource":    nil,
		"gitlab_naughty_resource": nil, // want `Resource "gitlab_naughty_resource" is missing a docs page named "naughty_resource.md"`
	},
}
