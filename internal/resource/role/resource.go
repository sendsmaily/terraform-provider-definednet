package role

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// NewResource creates a Defined.net Nebula host resource.
func NewResource() resource.Resource {
	return &Resource{}
}

// Resource is Defined.net Nebula host resource.
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
	resp.TypeName = fmt.Sprintf("%s_role", req.ProviderTypeName)
}

// Schema returns the resource's configuration schema.
func (r *Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = Schema
}

// Create creates Nebula hosts on Defined.net control plane.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state State

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := definednet.CreateRole(ctx, r.client, definednet.CreateRoleRequest{
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueString(),
		FirewallRules: lo.Map(state.FirewallRules, func(rule FirewallRule, _ int) (out definednet.FirewallRule) {
			out.Protocol = rule.Protocol.ValueString()
			out.Description = rule.Description.ValueString()
			out.AllowedRoleID = rule.AllowedRoleID.ValueString()

			if lo.IsNotNil(rule.PortRange) {
				out.PortRange = &definednet.PortRange{
					From: int(rule.PortRange.From.ValueInt32()),
					To:   int(rule.PortRange.To.ValueInt32()),
				}
			} else if !rule.Port.IsNull() {
				out.PortRange = &definednet.PortRange{
					From: int(rule.Port.ValueInt32()),
					To:   int(rule.Port.ValueInt32()),
				}
			}

			resp.Diagnostics.Append(rule.AllowedTags.ElementsAs(ctx, &out.AllowedTags, false)...)

			return out
		}),
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.Apply(ctx, role)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "created Defined.net role", map[string]any{
		"id":   state.ID.String(),
		"name": state.Name.String(),
	})
}

// Delete deletes Nebula hosts from Defined.net control plane.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state State

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := definednet.DeleteRole(ctx, r.client, definednet.DeleteRoleRequest{
		ID: state.ID.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
	}
}

// Read reads Nebula hosts from Defined.net control plane.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state State

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := definednet.GetRole(ctx, r.client, definednet.GetRoleRequest{
		ID: state.ID.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.Apply(ctx, role)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "refreshed Defined.net role", map[string]any{
		"id":   state.ID.String(),
		"name": state.Name.String(),
	})
}

// Update updates Nebula hosts on Defined.net control plane.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state State

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	role, err := definednet.UpdateRole(ctx, r.client, definednet.UpdateRoleRequest{
		ID:          state.ID.ValueString(),
		Name:        state.Name.ValueString(),
		Description: state.Description.ValueString(),
		FirewallRules: lo.Map(state.FirewallRules, func(rule FirewallRule, _ int) (out definednet.FirewallRule) {
			out.Protocol = rule.Protocol.ValueString()
			out.Description = rule.Description.ValueString()
			out.AllowedRoleID = rule.AllowedRoleID.ValueString()

			if lo.IsNotNil(rule.PortRange) {
				out.PortRange = &definednet.PortRange{
					From: int(rule.PortRange.From.ValueInt32()),
					To:   int(rule.PortRange.To.ValueInt32()),
				}
			} else if !rule.Port.IsNull() {
				out.PortRange = &definednet.PortRange{
					From: int(rule.Port.ValueInt32()),
					To:   int(rule.Port.ValueInt32()),
				}
			}

			resp.Diagnostics.Append(rule.AllowedTags.ElementsAs(ctx, &out.AllowedTags, false)...)

			return out
		}),
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	resp.Diagnostics.Append(state.Apply(ctx, role)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "updated Defined.net role", map[string]any{
		"id":   state.ID.String(),
		"name": state.Name.String(),
	})
}

// ImportState imports Nebula hosts from Defined.net control plane.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	role, err := definednet.GetRole(ctx, r.client, definednet.GetRoleRequest{
		ID: req.ID,
	})

	if err != nil {
		resp.Diagnostics.AddError("Request Failure", err.Error())
		return
	}

	var state State
	resp.Diagnostics.Append(state.Apply(ctx, role)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	tflog.Trace(ctx, "imported Defined.net role", map[string]any{
		"id":   state.ID.String(),
		"name": state.Name.String(),
	})
}
