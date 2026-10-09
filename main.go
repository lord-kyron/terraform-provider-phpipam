package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/lord-kyron/terraform-provider-phpipam/plugin/providers/phpipam"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{
		Debug:        debug,
		ProviderFunc: phpipam.Provider,
	}

	plugin.Serve(opts)
}
