package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/pavel-z1/terraform-provider-phpipam/plugin/providers/phpipam"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: phpipam.Provider,
	})
}
