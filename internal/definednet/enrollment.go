package definednet

import (
	"context"
	"net/http"
)

type (
	// Enrollment is a data model for Defined.net host enrollment.
	Enrollment struct {
		Host           Host           `json:"host"`
		EnrollmentCode EnrollmentCode `json:"enrollmentCode"`
	}

	// EnrollmentCode is a data model for Defined.net host enrollment code.
	EnrollmentCode struct {
		Code            string `json:"code"`
		LifetimeSeconds int    `json:"lifetimeSeconds"`
	}
)

// CreateEnrollment creates a Defined.net host enrollment.
func CreateEnrollment(ctx context.Context, client Client, req CreateEnrollmentRequest) (*Enrollment, error) {
	var resp Response[Enrollment]
	if err := client.Do(ctx, http.MethodPost, []string{"v1", "host-and-enrollment-code"}, req, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// CreateEnrollmentRequest is a request data model for CreateEnrollment endpoint.
type CreateEnrollmentRequest struct {
	NetworkID       string           `json:"networkID"`
	RoleID          string           `json:"roleID,omitempty"`
	Name            string           `json:"name"`
	StaticAddresses []string         `json:"staticAddresses"`
	ListenPort      int              `json:"listenPort"`
	IsLighthouse    bool             `json:"isLighthouse"`
	IsRelay         bool             `json:"isRelay"`
	Tags            []string         `json:"tags"`
	ConfigOverrides []ConfigOverride `json:"configOverrides"`
}
