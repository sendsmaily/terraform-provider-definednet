package host

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// State is the host resource's state.
type State struct {
	ID             types.String `tfsdk:"id"`
	NetworkID      types.String `tfsdk:"network_id"`
	RoleID         types.String `tfsdk:"role_id"`
	Name           types.String `tfsdk:"name"`
	IPAddress      types.String `tfsdk:"ip_address"`
	Tags           types.List   `tfsdk:"tags"`
	EnrollmentCode types.String `tfsdk:"enrollment_code"`
}

// ApplyHost applies Defined.net host information to the state.
func (s *State) ApplyHost(ctx context.Context, host *definednet.Host) (diags diag.Diagnostics) {
	s.ID = types.StringValue(host.ID)
	s.IPAddress = types.StringValue(host.IPAddress)
	s.Name = types.StringValue(host.Name)
	s.NetworkID = types.StringValue(host.NetworkID)
	s.RoleID = types.StringValue(host.RoleID)

	s.Tags, diags = types.ListValueFrom(ctx, types.StringType, host.Tags)

	return diags
}

// ApplyEnrollmentCode applies Defined.net enrollment code information to the state.
func (s *State) ApplyEnrollmentcode(_ context.Context, code *definednet.EnrollmentCode) (diags diag.Diagnostics) {
	s.EnrollmentCode = types.StringValue(code.Code)
	return diags
}
