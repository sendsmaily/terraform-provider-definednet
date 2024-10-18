package definednet

import (
	"context"
	"net/http"
)

// CreateEnrollmentCode creates a Defined.net host enrollment code.
func CreateEnrollmentCode(ctx context.Context, client Client, req CreateEnrollmentCodeRequest) (*CreateEnrollmentCodeResponse, error) {
	var resp CreateEnrollmentCodeResponse
	if err := client.Do(ctx, http.MethodPost, []string{"v1", "hosts", req.ID, "enrollment-code"}, nil, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

type (
	// CreateEnrollmentCodeRequest is a request data model for CreateEnrollmentCode endpoint.
	CreateEnrollmentCodeRequest struct {
		ID string
	}

	// CreateEnrollmentCodeResponse is a response data model for CreateEnrollmentCode endpoint.
	CreateEnrollmentCodeResponse struct {
		Data struct {
			Code            string `json:"code"`
			LifetimeSeconds int    `json:"lifetimeSeconds"`
		} `json:"data"`
	}
)
