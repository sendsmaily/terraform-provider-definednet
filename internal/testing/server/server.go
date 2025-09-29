package server

import (
	"log"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/onsi/ginkgo/v2"
	"github.com/samber/lo"
	"github.com/sendsmaily/terraform-provider-definednet/internal/definednet"
)

// New creates a fake Defined.net HTTP API server.
func New() *Server {
	middleware.DefaultLogger = middleware.RequestLogger(&middleware.DefaultLogFormatter{
		Logger: log.New(ginkgo.GinkgoWriter, "", log.LstdFlags),
	})

	mux := chi.NewMux()
	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	srv := &Server{
		Hosts: NewRepository[Host](),
		Roles: NewRepository[Role](),
	}

	// Hosts.
	mux.Post("/v1/host-and-enrollment-code", srv.createEnrollment)
	mux.Delete("/v1/hosts/{id}", srv.deleteHost)
	mux.Get("/v1/hosts/{id}", srv.getHost)
	mux.Put("/v2/hosts/{id}", srv.updateHost)

	// Roles.
	mux.Post("/v1/roles", srv.createRole)
	mux.Delete("/v1/roles/{id}", srv.deleteRole)
	mux.Get("/v1/roles/{id}", srv.getRole)
	mux.Put("/v1/roles/{id}", srv.updateRole)

	srv.server = httptest.NewServer(mux)

	return srv
}

// Server is a fake Defined.net HTTP API server.
type Server struct {
	Hosts *Repository[Host]
	Roles *Repository[Role]

	server *httptest.Server
}

// Close the fake HTTP API server.
func (s *Server) Close() {
	s.server.Close()
}

// Client returns a client for the fake HTTP API server.
func (s *Server) Client() definednet.Client {
	return lo.Must(definednet.NewClient(s.server.URL, "supersecret"))
}
