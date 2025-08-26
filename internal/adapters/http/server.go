package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/Ontair/dining-room/internal/core/ports"
)

const ReadHeaderTimeout = 10 * time.Second

type Server struct {
	http    *http.Server
	handler *DishesHandler
}

func NewServer(addr string, service ports.DishesService, logger *slog.Logger) *Server {
	handler := NewDishesHandler(service, logger)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /dish", handler.GetDishes)
	mux.HandleFunc("POST /dish", handler.CreateDish)

	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: ReadHeaderTimeout,
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
