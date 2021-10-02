package main

import (
	"github.com/gitlabhq/terraform-provider-gitlab/gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitlab.Provider})
}
