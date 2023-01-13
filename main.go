package main

import (
	"context"
	"terraform-provider-workos/workos"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.Background(), workos.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/vellanci/workos",
	})
}
