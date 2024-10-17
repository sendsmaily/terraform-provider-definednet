package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

const (
	// DefinednetApiEndpoint declares the default Defined.net HTTP API endpoint.
	DefinednetAPIEndpoint = "https://api.defined.net/"
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

var _ provider.Provider = (*Provider)(nil)

// Configuration declares the provider's configuration options.
type Configuration struct {
	Token types.String `tfsdk:"token"`
}

// Metadata returns the provider's metadata.
func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "definednet"
	resp.Version = p.version
}

// Schema returns the provider's configuration schema.
func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": schema.StringAttribute{
				Description: "Defined.net HTTP API token",
				Required:    true,
				Sensitive:   true,
			},
		},
	}
}

// Configure configures the provider with user passed options.
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Configuration

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client := definednet.NewClient(DefinednetAPIEndpoint, config.Token.String())
	resp.ResourceData = client
	resp.DataSourceData = client
}

// Resources returns a slice of resources available on the provider.
func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return nil
}

// DataSources returns a slice of data sources available on the provider.
func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
