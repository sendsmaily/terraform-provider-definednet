package role

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// State is the role resource's state.
type State struct {
	ID            types.String   `tfsdk:"id"`
	Name          types.String   `tfsdk:"name"`
	Description   types.String   `tfsdk:"description"`
	FirewallRules []FirewallRule `tfsdk:"rule"`
}

// FirewallRule is the role's firewall rule state.
type FirewallRule struct {
	Protocol      types.String       `tfsdk:"protocol"`
	Description   types.String       `tfsdk:"description"`
	AllowedRoleID types.String       `tfsdk:"allowed_role_id"`
	AllowedTags   types.Set          `tfsdk:"allowed_tags"`
	Port          types.Int32        `tfsdk:"port"`
	PortRange     *FirewallPortRange `tfsdk:"port_range"`
}

// FirewallPortRange is the firewall rule's port range state.
type FirewallPortRange struct {
	From types.Int32 `tfsdk:"from"`
	To   types.Int32 `tfsdk:"to"`
}

// Apply applies Defined.net role information to the state.
func (s *State) Apply(ctx context.Context, role *definednet.Role) (diags diag.Diagnostics) {
	s.ID = types.StringValue(role.ID)
	s.Name = types.StringValue(role.Name)
	s.Description = lo.If(lo.IsEmpty(role.Description), types.StringNull()).Else(types.StringValue(role.Description))
	s.FirewallRules = lo.Map(role.FirewallRules, func(rule definednet.FirewallRule, _ int) (r FirewallRule) {
		r.Protocol = types.StringValue(rule.Protocol)
		r.Description = lo.If(lo.IsNotEmpty(rule.Description), types.StringValue(rule.Description)).Else(types.StringNull())
		r.AllowedRoleID = lo.If(lo.IsNotEmpty(rule.AllowedRoleID), types.StringValue(rule.AllowedRoleID)).Else(types.StringNull())

		r.AllowedTags = func() basetypes.SetValue {
			tags, d := types.SetValueFrom(ctx, types.StringType, rule.AllowedTags)
			diags = append(diags, d...)
			return tags
		}()

		r.Port = func() types.Int32 {
			if lo.IsNil(rule.PortRange) || rule.PortRange.From != rule.PortRange.To {
				return types.Int32Null()
			}

			return types.Int32Value(int32(rule.PortRange.From))
		}()

		r.PortRange = func() *FirewallPortRange {
			if lo.IsNil(rule.PortRange) || rule.PortRange.From == rule.PortRange.To {
				return nil
			}

			return &FirewallPortRange{
				From: types.Int32Value(int32(rule.PortRange.From)),
				To:   types.Int32Value(int32(rule.PortRange.To)),
			}
		}()

		return r
	})

	return diags
}
