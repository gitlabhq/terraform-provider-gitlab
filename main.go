package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-gitlab/gitlab"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gitlab.Provider})
}
