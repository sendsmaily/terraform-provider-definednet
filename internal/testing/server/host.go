package server

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
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
