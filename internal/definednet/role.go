package definednet

import (
	"context"
	"net/http"
)

type (
	// Role is a data model for Defined.net role.
	Role struct {
		ID            string         `json:"id"`
		Name          string         `json:"name"`
		Description   string         `json:"description"`
		FirewallRules []FirewallRule `json:"firewallRules"`
	}

	// FirewallRule is a data model for Defined.net role's firewall rule.
	FirewallRule struct {
		Protocol      string     `json:"protocol"`
		Description   string     `json:"description"`
		AllowedRoleID string     `json:"allowedRoleID,omitempty"`
		AllowedTags   []string   `json:"allowedTags,omitempty"`
		PortRange     *PortRange `json:"portRange,omitempty"`
	}

	// PortRange is a data model for Defined.net role firewall rule's port range.
	PortRange struct {
		From int `json:"from"`
		To   int `json:"to"`
	}
)

// CreateRole creates a Defined.net role.
func CreateRole(ctx context.Context, client Client, req CreateRoleRequest) (*Role, error) {
	var resp Response[Role]
	if err := client.Do(ctx, http.MethodPost, []string{"v1", "roles"}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateRoleRequest is a request data model for CreateRole endpoint.
type CreateRoleRequest struct {
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	FirewallRules []FirewallRule `json:"firewallRules"`
}

// DeleteRole deletes a Defined.net role.
func DeleteRole(ctx context.Context, client Client, req DeleteRoleRequest) error {
	return client.Do(ctx, http.MethodDelete, []string{"v1", "roles", req.ID}, nil, nil)
}

// DeleteRoleRequest is a request data model for DeleteRole endpoint.
type DeleteRoleRequest struct {
	ID string
}

// GetRole retrieves a Defined.net role.
func GetRole(ctx context.Context, client Client, req GetRoleRequest) (*Role, error) {
	var resp Response[Role]
	if err := client.Do(ctx, http.MethodGet, []string{"v1", "roles", req.ID}, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetRoleRequest is a request data model for GetRole endpoint.
type GetRoleRequest struct {
	ID string
}

// UpdateRole updates a Defined.net role.
func UpdateRole(ctx context.Context, client Client, req UpdateRoleRequest) (*Role, error) {
	var resp Response[Role]
	if err := client.Do(ctx, http.MethodPut, []string{"v1", "roles", req.ID}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// UpdateRoleRequest is a request data model for UpdateRole endpoint.
type UpdateRoleRequest struct {
	ID            string         `json:"-"`
	Name          string         `json:"name"`
	Description   string         `json:"description"`
	FirewallRules []FirewallRule `json:"firewallRules"`
}
