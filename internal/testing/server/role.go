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

// Role is a data model for a Defined.net role.
type Role definednet.Role

// Key returns the role's repository key.
func (r Role) Key() string {
	return r.ID
}

func (s *Server) createRole(w http.ResponseWriter, r *http.Request) {
	var req definednet.CreateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	state := Role{
		ID:            fmt.Sprintf("role-%s", strings.ToUpper(lo.RandomString(8, lo.AlphanumericCharset))),
		Name:          req.Name,
		Description:   req.Description,
		FirewallRules: req.FirewallRules,
	}

	if err := s.Roles.Add(state); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Role]{
		Data: definednet.Role(state),
	}); err != nil {
		panic(err)
	}
}

func (s *Server) getRole(w http.ResponseWriter, r *http.Request) {
	state, err := s.Roles.Get(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Role]{
		Data: definednet.Role(*state),
	}); err != nil {
		panic(err)
	}
}

func (s *Server) updateRole(w http.ResponseWriter, r *http.Request) {
	var req definednet.UpdateRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		panic(err)
	}

	state, err := s.Roles.Get(chi.URLParam(r, "id"))
	if err != nil {
		panic(err)
	}

	state.Name = req.Name
	state.Description = req.Description
	state.FirewallRules = req.FirewallRules

	if err := s.Roles.Replace(*state); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(definednet.Response[definednet.Role]{
		Data: definednet.Role(*state),
	}); err != nil {
		panic(err)
	}
}

func (s *Server) deleteRole(w http.ResponseWriter, r *http.Request) {
	if err := s.Roles.Remove(chi.URLParam(r, "id")); err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusOK)
}
