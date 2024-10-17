package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// New creates a Defined.net Terraform provider.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			version: version,
		}
	}
}

// Provider is the Defined.net Terraform provider.
type Provider struct {
	version string
}

var _ provider.Provider = &Provider{}

// Configuration declares the provider's configuration options.
type Configuration struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

// Metadata returns the provider's metadata.
func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "definednet"
	resp.Version = p.version
}

// Schema returns the provider's configuration schema.
func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}

// Configure configures the provider with user passed options.
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

// Resources returns a slice of resources available on the provider.
func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

// DataSources returns a slice of data sources available on the provider.
func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
