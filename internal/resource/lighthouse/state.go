package lighthouse

import (
	"context"
	"fmt"
	"net"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// State is the lighthouse resource's state.
type State struct {
	ID              types.String `tfsdk:"id"`
	NetworkID       types.String `tfsdk:"network_id"`
	RoleID          types.String `tfsdk:"role_id"`
	StaticAddresses types.List   `tfsdk:"static_addresses"`
	ListenPort      types.Int32  `tfsdk:"listen_port"`
	Name            types.String `tfsdk:"name"`
	IPAddress       types.String `tfsdk:"ip_address"`
	Tags            types.List   `tfsdk:"tags"`
	EnrollmentCode  types.String `tfsdk:"enrollment_code"`
}

// ApplyEnrollment applies Defined.net host enrollment information to the state.
func (s *State) ApplyEnrollment(ctx context.Context, enrollment *definednet.Enrollment) (diags diag.Diagnostics) {
	diags.Append(s.ApplyHost(ctx, &enrollment.Host)...)
	s.EnrollmentCode = types.StringValue(enrollment.EnrollmentCode.Code)

	return diags
}

// ApplyHost applies Defined.net lighthouse information to the state.
func (s *State) ApplyHost(ctx context.Context, lighthouse *definednet.Host) (diags diag.Diagnostics) {
	staticAddrs, d := types.ListValueFrom(ctx, types.StringType, lo.Map(lighthouse.StaticAddresses, func(addr string, idx int) string {
		a, err := net.ResolveTCPAddr("tcp", addr)
		if err != nil {
			// TODO: what is the correct path to return here? The state? The HTTP API response?
			diags.AddAttributeError(
				path.Empty().AtMapKey("static_addresses").AtListIndex(idx),
				fmt.Sprintf("Invalid Address %q", addr),
				err.Error(),
			)

			return "<nil>"
		}

		return a.IP.String()
	}))

	diags.Append(d...)

	s.ID = types.StringValue(lighthouse.ID)
	s.NetworkID = types.StringValue(lighthouse.NetworkID)
	s.StaticAddresses = staticAddrs
	s.ListenPort = types.Int32Value(int32(lighthouse.ListenPort))
	s.Name = types.StringValue(lighthouse.Name)
	s.IPAddress = types.StringValue(lighthouse.IPAddress)

	if !lo.IsEmpty(lighthouse.RoleID) {
		s.RoleID = types.StringValue(lighthouse.RoleID)
	}

	if len(lighthouse.Tags) > 0 {
		tags, d := types.ListValueFrom(ctx, types.StringType, lighthouse.Tags)
		s.Tags = tags
		diags.Append(d...)
	}

	return diags
}
