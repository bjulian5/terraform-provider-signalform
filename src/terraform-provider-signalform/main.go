package main

import (
	"github.com/Yelp/terraform-provider-signalform/src/terraform-provider-signalform/signalform"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: signalform.Provider,
	})
}
