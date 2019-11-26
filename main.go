package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/lord-kyron/terraform-provider-phpipam-0.3.0/plugin/providers/phpipam"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: phpipam.Provider,
	})
}
