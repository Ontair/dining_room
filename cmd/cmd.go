package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// "github.com/Ontair/dining-room/internal/".
	httpAdapter "github.com/Ontair/dining-room/internal/adapters/http"
	"github.com/Ontair/dining-room/internal/core/memory"
	"github.com/Ontair/dining-room/internal/core/service"
)

var shutdownTimeout = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}
	l.Info("сервер запущен", slog.String("port", addr))

	repo := memory.NewMemoryDishesRepository()
	dininService := service.NewDishesService(repo, l)
	server := httpAdapter.NewServer(addr, dininService, l)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error("failed to start server: %v", slog.Any("error", err))
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	l.Info("Сигнал закрытия")

	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Info("server forced to shutdown: %v", slog.Any("error", err))
	}

	l.Info("сервер закрыт")
}
