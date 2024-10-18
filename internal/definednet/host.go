package definednet

import (
	"context"
	"net/http"
)

// CreateHost creates a Defined.net host.
func CreateHost(ctx context.Context, client Client, req CreateHostRequest) (*CreateHostResponse, error) {
	var resp Response[CreateHostResponse]
	if err := client.Do(ctx, http.MethodPost, []string{"v1", "hosts"}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

type (
	// CreateHostRequest is a request data model for CreateHost endpoint.
	CreateHostRequest struct {
		NetworkID       string   `json:"networkID"`
		RoleID          string   `json:"roleID"`
		Name            string   `json:"name"`
		IPAddress       string   `json:"ipAddress"`
		StaticAddresses []string `json:"staticAddresses"`
		ListenPort      int      `json:"listenPort"`
		IsLighthouse    bool     `json:"isLighthouse"`
		IsRelay         bool     `json:"isRelay"`
		Tags            []string `json:"tags"`
	}

	// CreateHostResponse is a response data model for CreateHost endpoint.
	CreateHostResponse struct {
		ID              string   `json:"id"`
		NetworkID       string   `json:"networkID"`
		RoleID          string   `json:"roleID"`
		Name            string   `json:"name"`
		IPAddress       string   `json:"ipAddress"`
		StaticAddresses []string `json:"staticAddresses"`
		ListenPort      int      `json:"listenPort"`
		IsLighthouse    bool     `json:"isLighthouse"`
		IsRelay         bool     `json:"isRelay"`
		Tags            []string `json:"tags"`
	}
)

// DeleteHost deletes a Defined.net host.
func DeleteHost(ctx context.Context, client Client, req DeleteHostRequest) error {
	return client.Do(ctx, http.MethodDelete, []string{"v1", "hosts", req.ID}, nil, nil)
}

type (
	// DeleteHostRequest is a request data model for DeleteHost endpoint.
	DeleteHostRequest struct {
		ID string
	}
)

// GetHost retrieves a Defined.net host.
func GetHost(ctx context.Context, client Client, req GetHostRequest) (*GetHostResponse, error) {
	var resp Response[GetHostResponse]
	if err := client.Do(ctx, http.MethodGet, []string{"v1", "hosts", req.ID}, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

type (
	// GetHostRequest is a request data model for GetHost endpoint.
	GetHostRequest struct {
		ID string
	}

	// GetHostResponse is a response data model for GetHost endpoint.
	GetHostResponse struct {
		ID              string   `json:"id"`
		NetworkID       string   `json:"networkID"`
		RoleID          string   `json:"roleID"`
		Name            string   `json:"name"`
		IPAddress       string   `json:"ipAddress"`
		StaticAddresses []string `json:"staticAddresses"`
		ListenPort      int      `json:"listenPort"`
		IsLighthouse    bool     `json:"isLighthouse"`
		IsRelay         bool     `json:"isRelay"`
		Tags            []string `json:"tags"`
	}
)

// UpdateHost updates a Defined.net host.
func UpdateHost(ctx context.Context, client Client, req UpdateHostRequest) (*UpdateHostResponse, error) {
	var resp Response[UpdateHostResponse]
	if err := client.Do(ctx, http.MethodPut, []string{"v1", "hosts", req.ID}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

type (
	// UpdateHostRequest is a request data model for UpdateHost endpoint.
	UpdateHostRequest struct {
		ID              string   `json:"-"`
		RoleID          string   `json:"roleID"`
		Name            string   `json:"name"`
		StaticAddresses []string `json:"staticAddresses"`
		ListenPort      int      `json:"listenPort"`
		Tags            []string `json:"tags"`
	}

	// UpdateHostResponse is a response data model for UpdateHost endpoint.
	UpdateHostResponse struct {
		ID              string   `json:"id"`
		NetworkID       string   `json:"networkID"`
		RoleID          string   `json:"roleID"`
		Name            string   `json:"name"`
		IPAddress       string   `json:"ipAddress"`
		StaticAddresses []string `json:"staticAddresses"`
		ListenPort      int      `json:"listenPort"`
		IsLighthouse    bool     `json:"isLighthouse"`
		IsRelay         bool     `json:"isRelay"`
		Tags            []string `json:"tags"`
	}
)
