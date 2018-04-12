package main

import (
	"github.com/claranet/terraform-provider-gitlab/gitlab"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitlab.Provider})
}
