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

	counter := NewRequestCounter("get_count.txt", "post_count.txt")

	// Define our middleware stack
	// These run in the order given
	stackDish := []Middleware{
		CreateCountAndWriteRequestMiddleware(counter),
	}

	mux.HandleFunc("GET /dish", CompileMiddleware(handler.GetDishes, stackDish))
	mux.HandleFunc("POST /dish", CompileMiddleware(handler.CreateDish, stackDish))

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
