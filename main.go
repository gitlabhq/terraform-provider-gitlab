package main

import (
	"github.com/gitlabhq/terraform-provider-gitlab/gitlab"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitlab.Provider})
}
