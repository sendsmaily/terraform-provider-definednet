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
	Metrics         *Metrics     `tfsdk:"metrics"`
}

// Metrics is the host metrics exporter's state.
type Metrics struct {
	Enabled            types.Bool   `tfsdk:"enabled"`
	Listen             types.String `tfsdk:"listen"`
	Path               types.String `tfsdk:"path"`
	Namespace          types.String `tfsdk:"namespace"`
	Subsystem          types.String `tfsdk:"subsystem"`
	EnableExtraMetrics types.Bool   `tfsdk:"enable_extra_metrics"`
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

	tags := types.ListNull(types.StringType)
	if len(lighthouse.Tags) > 0 {
		tags, d = types.ListValueFrom(ctx, types.StringType, lighthouse.Tags)
		diags.Append(d...)
	}

	s.ID = types.StringValue(lighthouse.ID)
	s.Name = types.StringValue(lighthouse.Name)
	s.NetworkID = types.StringValue(lighthouse.NetworkID)
	s.StaticAddresses = staticAddrs
	s.ListenPort = types.Int32Value(int32(lighthouse.ListenPort))
	s.IPAddress = types.StringValue(lighthouse.IPAddress)
	s.Tags = tags

	s.RoleID = types.StringNull()
	if lo.IsNotEmpty(lighthouse.RoleID) {
		s.RoleID = types.StringValue(lighthouse.RoleID)
	}

	metricsConfig := lo.Reduce(lighthouse.ConfigOverrides, func(m Metrics, o definednet.ConfigOverride, _ int) Metrics {
		switch o.Key {
		case "stats.type":
			if o.Value != "prometheus" {
				diags.AddError("Unsupported Metrics Backend", fmt.Sprintf("Expected 'prometheus', got '%s'", o.Value))
			}

			m.Enabled = types.BoolValue(true)

		case "stats.listen":
			v, err := convert[string](o.Value)
			if err != nil {
				diags.AddAttributeError(path.Root("metrics").AtMapKey("listen"), "Invalid Value", err.Error())
			}

			m.Listen = types.StringValue(v)

		case "stats.path":
			v, err := convert[string](o.Value)
			if err != nil {
				diags.AddAttributeError(path.Root("metrics").AtMapKey("path"), "Invalid Value", err.Error())
			}

			m.Path = types.StringValue(v)

		case "stats.namespace":
			v, err := convert[string](o.Value)
			if err != nil {
				diags.AddAttributeError(path.Root("metrics").AtMapKey("namespace"), "Invalid Value", err.Error())
			}

			m.Namespace = types.StringValue(v)

		case "stats.subsystem":
			v, err := convert[string](o.Value)
			if err != nil {
				diags.AddAttributeError(path.Root("metrics").AtMapKey("subsystem"), "Invalid Value", err.Error())
			}

			m.Subsystem = types.StringValue(v)

		case "stats.lighthouse_metrics":
			v, err := convert[bool](o.Value)
			if err != nil {
				diags.AddAttributeError(path.Root("metrics").AtMapKey("enable_extra_metrics"), "Invalid Value", err.Error())
			}

			m.EnableExtraMetrics = types.BoolValue(v)
		}

		return m
	}, Metrics{})

	if lo.IsNotEmpty(metricsConfig) {
		s.Metrics = &metricsConfig
	}

	return diags
}

func convert[T any](val any) (T, error) {
	if val, ok := val.(T); ok {
		return val, nil
	}

	var t T
	return t, fmt.Errorf("unexpected type: wanted %T, got %T", t, val)
}
