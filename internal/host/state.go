package host

import "github.com/hashicorp/terraform-plugin-framework/types"

// State is the host resource's state.
type State struct {
	ID        types.String `tfsdk:"id"`
	Name      types.String `tfsdk:"name"`
	NetworkID types.String `tfsdk:"network_id"`
	RoleID    types.String `tfsdk:"role_id"`
	Tags      types.List   `tfsdk:"tags"`

	IPAddress      types.String `tfsdk:"ip_address"`
	EnrollmentCode types.String `tfsdk:"enrollment_code"`
}
