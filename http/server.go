package http

import (
	"net"
	"net/http"
)

// Server represents the server
type Server struct {
	ln net.Listener

	Handler *Handler

	Addr string
}

// DefaultAddr is the default bind address.
const DefaultAddr = ":3000"

// NewServer returns a new instance of Server.
func NewServer() *Server {
	return &Server{
		Addr: DefaultAddr,
	}
}

// Open is used to start the server
func (s *Server) Open() error {
	ln, err := net.Listen("tcp", s.Addr)

	if err != nil {
		return err
	}

	s.ln = ln

	// Start HTTP server.
	go func() { http.Serve(s.ln, s.Handler) }()

	return nil
}

// Close is used to shutdown the server
func (s *Server) Close() error {
	if s.ln != nil {
		s.ln.Close()
	}

	return nil
}
