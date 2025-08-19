package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	httpAdapter "github.com/Ontair/dining-room/internal/adapters/http"
	"github.com/Ontair/dining-room/internal/adapters/repository"

)

func main(){
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	addr := os.Getenv("ADDR")
	if addr == "" {
		addr = ":8080"
	}

	l := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := repository.NewMemoryDinnerRepository()
	dininService := service.NewDinnerService(repo, l)
	server := httpAdapter.NewServer(addr, dininService, l)

	go func(){
		if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed){
			l.Error("failed to start server: %v", err)
		}
	}()

	<-ctx.Done()
	
}