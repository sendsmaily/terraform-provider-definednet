package provider

import (
	"context"
	_ "embed"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
	"github.com/sendsmaily/terraform-provider-definednet/internal/resource/host"
	"github.com/sendsmaily/terraform-provider-definednet/internal/resource/lighthouse"
	"github.com/sendsmaily/terraform-provider-definednet/internal/resource/role"
)

const (
	// DefinednetApiEndpoint declares the default Defined.net HTTP API endpoint.
	DefinednetAPIEndpoint = "https://api.defined.net/"
)

// ClientFactory creates Defined.net clients.
type ClientFactory func(endpointURI string, token string, version string) (definednet.Client, error)

// New creates a Defined.net Terraform provider.
func New(clientFactory ClientFactory, version string) func() provider.Provider {
	return func() provider.Provider {
		return &Provider{
			clientFactory: clientFactory,
			version:       version,
		}
	}
}

// Provider is the Defined.net Terraform provider.
type Provider struct {
	clientFactory ClientFactory
	version       string
}

// Configuration declares the provider's configuration options.
type Configuration struct {
	Token types.String `tfsdk:"token"`
}

var _ provider.Provider = (*Provider)(nil)

// Metadata returns the provider's metadata.
func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "definednet"
	resp.Version = p.version
}

// Schema returns the provider's configuration schema.
func (p *Provider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = Schema
}

// Configure configures the provider with user passed options.
func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config Configuration

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	client, err := p.clientFactory(DefinednetAPIEndpoint, config.Token.ValueString(), p.version)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Configuration", err.Error())
		return
	}

	resp.ResourceData = client
	resp.DataSourceData = client
}

// Resources returns a slice of resources available on the provider.
func (p *Provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		lighthouse.NewResource,
		host.NewResource,
		role.NewResource,
	}
}

// DataSources returns a slice of data sources available on the provider.
func (p *Provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}
