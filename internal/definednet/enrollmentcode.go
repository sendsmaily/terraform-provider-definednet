package definednet

import (
	"context"
	"net/http"
)

// EnrollmentCode is a data model for Defined.net host enrollment code.
type EnrollmentCode struct {
	Code            string `json:"code"`
	LifetimeSeconds int    `json:"lifetimeSeconds"`
}

// CreateEnrollmentCode creates a Defined.net host enrollment code.
func CreateEnrollmentCode(ctx context.Context, client Client, req CreateEnrollmentCodeRequest) (*EnrollmentCode, error) {
	var resp Response[EnrollmentCode]
	if err := client.Do(ctx, http.MethodPost, []string{"v1", "hosts", req.ID, "enrollment-code"}, nil, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateEnrollmentCodeRequest is a request data model for CreateEnrollmentCode endpoint.
type CreateEnrollmentCodeRequest struct {
	ID string
}
