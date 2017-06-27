package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/sl1pm4t/terraform-provider-gobetween/gobetween"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: gobetween.Provider,
	})
}
