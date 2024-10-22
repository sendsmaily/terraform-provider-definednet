package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// Host is a data model for a Defined.net host.
type Host struct {
	Host           definednet.Host
	EnrollmentCode definednet.EnrollmentCode
}

// Key returns the host's repository key.
func (h Host) Key() string {
	return h.Host.ID
}

func (s *Server) createHost(w http.ResponseWriter, r *http.Request) {
	var req definednet.CreateHostRequest
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
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Host]{
		Data: state.Host,
	}); err != nil {
		panic(err)
	}
}

func (s *Server) getHost(w http.ResponseWriter, r *http.Request) {
	state, err := s.Hosts.Get(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Host]{
		Data: state.Host,
	}); err != nil {
		panic(err)
	}
}

func (s *Server) updateHost(w http.ResponseWriter, r *http.Request) {
	var req definednet.UpdateHostRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	state, err := s.Hosts.Get(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	state.Host.Name = req.Name
	state.Host.RoleID = req.RoleID
	state.Host.StaticAddresses = req.StaticAddresses
	state.Host.ListenPort = req.ListenPort
	state.Host.Tags = req.Tags

	if err := s.Hosts.Replace(*state); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Host]{
		Data: state.Host,
	}); err != nil {
		panic(err)
	}
}

func (s *Server) deleteHost(w http.ResponseWriter, r *http.Request) {
	if err := s.Hosts.Remove(chi.URLParam(r, "id")); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) createEnrollmentCode(w http.ResponseWriter, r *http.Request) {
	state, err := s.Hosts.Get(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	state.EnrollmentCode = definednet.EnrollmentCode{
		Code:            lo.RandomString(16, lo.AlphanumericCharset),
		LifetimeSeconds: 300,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.EnrollmentCode]{
		Data: state.EnrollmentCode,
	}); err != nil {
		panic(err)
	}
}
