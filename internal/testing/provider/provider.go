package provider

import (
	"context"

	tfprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	tfresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
	"github.com/sendsmaily/terraform-provider-definednet/internal/provider"
)

// New creates a fake provider.
func New(client definednet.Client) *Provider {
	return &Provider{
		Provider: provider.New("fake")(),
		client:   client,
	}
}

// Provider is a fake provider.
type Provider struct {
	tfprovider.Provider

	client    definednet.Client
	resources []tfresource.Resource
}

var _ tfprovider.Provider = (*Provider)(nil)

// WithResource configures the fake provider with a managed resource.
func (p Provider) WithResource(res tfresource.Resource) *Provider {
	p.resources = append(p.resources, res)
	return &p
}

// Configure configures the fake provider with user passed options.
func (p *Provider) Configure(ctx context.Context, req tfprovider.ConfigureRequest, resp *tfprovider.ConfigureResponse) {
	resp.ResourceData = p.client
	resp.DataSourceData = p.client
}

// Resources returns a slice of resources available on the fake provider.
func (p *Provider) Resources(ctx context.Context) []func() tfresource.Resource {
	return lo.Map(p.resources, func(res tfresource.Resource, _ int) func() tfresource.Resource {
		return func() tfresource.Resource {
			return res
		}
	})
}
