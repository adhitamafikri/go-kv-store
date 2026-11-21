package server

import (
	"context"
	"errors"
	"net"

	kv "github.com/adhitamafikri/go-kv-store"
)

// Used in Stage 9
type Server struct{}

func NewServer(db *kv.DB) *Server
func (s *Server) Serve(l net.Listener) error

// Used inStage 10
var ErrServerClosed = errors.New("kv: Server closed")

func (s *Server) Shutdown(ctx context.Context) error
