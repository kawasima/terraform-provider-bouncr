package main

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/kawasima/terraform-provider-bouncr/internal/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.New, providerserver.ServeOpts{
		Address: "registry.terraform.io/kawasima/bouncr",
	})
	if err != nil {
		log.Fatal(err)
	}
}
