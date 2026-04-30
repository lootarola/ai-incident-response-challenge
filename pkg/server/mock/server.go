package mock

import (
	"context"

	"github.com/lootarola/ai-incident-response-challenge/pkg/server"
)

func init() {
	server.Register("mock", newServer)
}

type Server struct{}

func newServer(_ server.Config) (server.Server, error) { return &Server{}, nil }
func (s *Server) Start(_ context.Context) error        { return nil }
