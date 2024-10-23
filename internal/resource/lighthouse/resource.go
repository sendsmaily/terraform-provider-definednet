package lighthouse

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// NewResource creates a Defined.net Nebula lighthouse resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Resource is Defined.net Nebula lighthouse resource.
type Resource struct {
	client definednet.Client
}

var _ resource.Resource = (*Resource)(nil)
var _ resource.ResourceWithConfigure = (*Resource)(nil)
var _ resource.ResourceWithImportState = (*Resource)(nil)

// Configure configures the resource.
func (r *Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(definednet.Client)
	if !ok {
		resp.Diagnostics.AddError("Invalid Configuration", "The provider specifies an invalid client type")
		return
	}

	r.client = client
}

// Metadata returns the resource's metadata.
func (r *Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = fmt.Sprintf("%s_lighthouse", req.ProviderTypeName)
}

// Schema returns the resource's configuration schema.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = Schema
}

// Create creates Nebula lighthouses on Defined.net control plane.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state State

	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var staticAddrs []string
	resp.Diagnostics.Append(state.StaticAddresses.ElementsAs(ctx, &staticAddrs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tags []string
	resp.Diagnostics.Append(state.Tags.ElementsAs(ctx, &tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := definednet.CreateHost(ctx, r.client, definednet.CreateHostRequest{
		NetworkID: state.NetworkID.ValueString(),
		RoleID:    state.RoleID.ValueString(),
		Name:      state.Name.ValueString(),
		StaticAddresses: lo.Map(staticAddrs, func(addr string, _ int) string {
			return fmt.Sprintf("%s:%d", addr, state.ListenPort.ValueInt32())
		}),
		ListenPort:   int(state.ListenPort.ValueInt32()),
		IsLighthouse: true,
		IsRelay:      false,
		Tags:         tags,
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.ApplyHost(ctx, host)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "created Defined.net lighthouse", map[string]any{
		"id":               state.ID.String(),
		"network_id":       state.NetworkID.String(),
		"role_id":          state.RoleID.String(),
		"name":             state.Name.String(),
		"listen_port":      state.ListenPort.String(),
		"static_addresses": state.StaticAddresses.String(),
		"tags":             state.Tags.String(),
	})

	code, err := definednet.CreateEnrollmentCode(ctx, r.client, definednet.CreateEnrollmentCodeRequest{
		ID: state.ID.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	tflog.Trace(ctx, "created Defined.net enrollment code", map[string]any{
		"id": state.ID.String(),
	})

	resp.Diagnostics.Append(state.ApplyEnrollmentcode(ctx, code)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes Nebula lighthouses from Defined.net control plane.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state State

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := definednet.DeleteHost(ctx, r.client, definednet.DeleteHostRequest{
		ID: state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
	}
}

// Read reads Nebula lighthouses from Defined.net control plane.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state State

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := definednet.GetHost(ctx, r.client, definednet.GetHostRequest{
		ID: state.ID.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.ApplyHost(ctx, host)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "refreshed Defined.net lighthouse", map[string]any{
		"id":               state.ID.String(),
		"network_id":       state.NetworkID.String(),
		"role_id":          state.RoleID.String(),
		"name":             state.Name.String(),
		"listen_port":      state.ListenPort.String(),
		"static_addresses": state.StaticAddresses.String(),
		"tags":             state.Tags.String(),
	})
}

// Update updates Nebula lighthouses on Defined.net control plane.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state State

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var staticAddrs []string
	resp.Diagnostics.Append(state.StaticAddresses.ElementsAs(ctx, &staticAddrs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var tags []string
	resp.Diagnostics.Append(state.Tags.ElementsAs(ctx, &tags, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	host, err := definednet.UpdateHost(ctx, r.client, definednet.UpdateHostRequest{
		ID:     state.ID.ValueString(),
		RoleID: state.RoleID.ValueString(),
		Name:   state.Name.ValueString(),
		StaticAddresses: lo.Map(staticAddrs, func(addr string, _ int) string {
			return fmt.Sprintf("%s:%d", addr, state.ListenPort.ValueInt32())
		}),
		ListenPort: int(state.ListenPort.ValueInt32()),
		Tags:       tags,
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.ApplyHost(ctx, host)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "updated Defined.net lighthouse", map[string]any{
		"id":               state.ID.String(),
		"network_id":       state.NetworkID.String(),
		"role_id":          state.RoleID.String(),
		"name":             state.Name.String(),
		"listen_port":      state.ListenPort.String(),
		"static_addresses": state.StaticAddresses.String(),
		"tags":             state.Tags.String(),
	})
}

// ImportState imports Nebula lighthouses from Defined.net control plane.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	host, err := definednet.GetHost(ctx, r.client, definednet.GetHostRequest{
		ID: req.ID,
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	var state State
	resp.Diagnostics.Append(state.ApplyHost(ctx, host)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(state.ApplyEnrollmentcode(ctx, &definednet.EnrollmentCode{})...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "imported Defined.net lighthouse", map[string]any{
		"id":               state.ID.String(),
		"network_id":       state.NetworkID.String(),
		"role_id":          state.RoleID.String(),
		"name":             state.Name.String(),
		"listen_port":      state.ListenPort.String(),
		"static_addresses": state.StaticAddresses.String(),
		"tags":             state.Tags.String(),
	})
}
