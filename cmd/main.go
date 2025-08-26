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

	httpAdapter "github.com/Ontair/dining-room/internal/adapters/http"
	"github.com/Ontair/dining-room/internal/adapters/memory"
	"github.com/Ontair/dining-room/internal/core/service"
)

const shutdownTimeout = 5 * time.Second

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	repo := memory.NewMemoryDishesRepository()
	dininService := service.NewDishesService(repo, l)
	server := httpAdapter.NewServer(addr, dininService, l)

	l.Info(
		"server is running",
		slog.String("port", server.Addr()),
	)

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			l.Error(
				"failed to start the server",
				slog.Any("error", err),
			)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	l.Info("closing signal")

	if err := server.Shutdown(shutdownCtx); err != nil {
		l.Info(
			"power server is turned off",
			slog.Any("error", err),
		)
	}

	l.Info("server closed")
}
