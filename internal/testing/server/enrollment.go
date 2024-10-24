package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

func (s *Server) createEnrollment(w http.ResponseWriter, r *http.Request) {
	var req definednet.CreateEnrollmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	state := Host{
		Host: definednet.Host{
			ID:        fmt.Sprintf("host-%s", strings.ToUpper(lo.RandomString(8, lo.AlphanumericCharset))),
			NetworkID: req.NetworkID,
			RoleID:    req.RoleID,
			Name:      req.Name,
			IPAddress: func() string {
				if !lo.IsEmpty(req.IPAddress) {
					return req.IPAddress
				}

				return "10.0.0.1"
			}(),
			StaticAddresses: req.StaticAddresses,
			ListenPort:      req.ListenPort,
			IsLighthouse:    req.IsLighthouse,
			IsRelay:         req.IsRelay,
			Tags:            req.Tags,
		},
	}

	if err := s.Hosts.Add(state); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Enrollment]{
		Data: definednet.Enrollment{
			Host: state.Host,
			EnrollmentCode: definednet.EnrollmentCode{
				Code:            lo.RandomString(32, lo.AllCharset),
				LifetimeSeconds: 300,
			},
		},
	}); err != nil {
		panic(err)
	}
}
