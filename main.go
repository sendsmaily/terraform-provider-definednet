// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
	"github.com/sendsmaily/terraform-provider-definednet/internal/provider"
)

// Configured by Goreleaser during build.
var (
	version string = "dev"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	if err := providerserver.Serve(
		context.Background(),
		provider.New(definednet.NewClient, version),
		providerserver.ServeOpts{
			Address: "registry.terraform.io/sendsmaily/definednet",
			Debug:   debug,
		},
	); err != nil {
		log.Fatal(err.Error())
	}
}
