package definednet

import (
	"context"
	"net/http"
)

// Host is a data model for Defined.net host.
type Host struct {
	ID              string           `json:"id"`
	NetworkID       string           `json:"networkID"`
	RoleID          string           `json:"roleID,omitempty"`
	Name            string           `json:"name"`
	IPAddress       string           `json:"ipAddress"`
	StaticAddresses []string         `json:"staticAddresses"`
	ListenPort      int              `json:"listenPort"`
	IsLighthouse    bool             `json:"isLighthouse"`
	IsRelay         bool             `json:"isRelay"`
	Tags            []string         `json:"tags"`
	ConfigOverrides []ConfigOverride `json:"configOverrides"`
}

// DeleteHost deletes a Defined.net host.
func DeleteHost(ctx context.Context, client Client, req DeleteHostRequest) error {
	return client.Do(ctx, http.MethodDelete, []string{"v1", "hosts", req.ID}, nil, nil)
}

// DeleteHostRequest is a request data model for DeleteHost endpoint.
type DeleteHostRequest struct {
	ID string
}

// GetHost retrieves a Defined.net host.
func GetHost(ctx context.Context, client Client, req GetHostRequest) (*Host, error) {
	var resp Response[Host]
	if err := client.Do(ctx, http.MethodGet, []string{"v1", "hosts", req.ID}, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetHostRequest is a request data model for GetHost endpoint.
type GetHostRequest struct {
	ID string
}

// UpdateHost updates a Defined.net host.
func UpdateHost(ctx context.Context, client Client, req UpdateHostRequest) (*Host, error) {
	var resp Response[Host]
	if err := client.Do(ctx, http.MethodPut, []string{"v2", "hosts", req.ID}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// UpdateHostRequest is a request data model for UpdateHost endpoint.
type UpdateHostRequest struct {
	ID              string           `json:"-"`
	RoleID          string           `json:"roleID,omitempty"`
	Name            string           `json:"name"`
	StaticAddresses []string         `json:"staticAddresses"`
	ListenPort      int              `json:"listenPort"`
	Tags            []string         `json:"tags"`
	ConfigOverrides []ConfigOverride `json:"configOverrides"`
}
