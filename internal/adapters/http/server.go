package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/Ontair/dining-room/internal/core/ports"
)

type Server struct {
	http    *http.Server
	handler *DiningHandler
}

func NewServer(addr string, service ports.DiningService, logger *slog.Logger) *Server {
	handler := NewDiningHandler(service, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dinings", handler.GetDinings)
	mux.HandleFunc("POST /dinings", handler.CreateDinings)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &Server{
		http:    httpServer,
		handler: handler,
	}
}

func (s *Server) ListenAndServe() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}

func (s *Server) Addr() string {
	return s.http.Addr
}
