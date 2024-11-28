package host

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/samber/lo"
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
	Metrics        *Metrics     `tfsdk:"metrics"`
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

// ApplyHost applies Defined.net host information to the state.
func (s *State) ApplyHost(ctx context.Context, host *definednet.Host) (diags diag.Diagnostics) {
	s.ID = types.StringValue(host.ID)
	s.IPAddress = types.StringValue(host.IPAddress)
	s.Name = types.StringValue(host.Name)
	s.NetworkID = types.StringValue(host.NetworkID)

	s.RoleID = types.StringNull()
	if lo.IsNotEmpty(host.RoleID) {
		s.RoleID = types.StringValue(host.RoleID)
	}

	s.Tags = types.ListNull(types.StringType)
	if len(host.Tags) > 0 {
		s.Tags, diags = types.ListValueFrom(ctx, types.StringType, host.Tags)
	}

	metricsConfig := lo.Reduce(host.ConfigOverrides, func(m Metrics, o definednet.ConfigOverride, _ int) Metrics {
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

		case "stats.message_metrics":
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
