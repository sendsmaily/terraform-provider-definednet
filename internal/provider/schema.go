package provider

import (
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
)

// Schema is the provider's configuration schema.
var Schema = schema.Schema{
	Description:         "Interact with Defined.net's managed Nebula control plane.",
	MarkdownDescription: providerDescription,
	Attributes: map[string]schema.Attribute{
		"token": schema.StringAttribute{
			Description: "Defined.net HTTP API token",
			Required:    true,
			Sensitive:   true,
		},
	},
}

//go:embed docs/provider.md
var providerDescription string
