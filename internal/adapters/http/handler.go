package http

import (
	"log/slog"
	"net/http"

	"github.com/Ontair/dining-room/internal/core/ports"
)

type DiningHandler struct{
	service ports.DiningService
	logger *slog.Logger
}

func NewDiningHandler(service ports.DiningService, logger *slog.Logger) *DiningHandler{
	return &DiningHandler{
		service: service,
		logger: logger,
	}
}

func (* DiningHandler) GetDinings(w http.ResponseWriter, r *http.Request){

}

func (* DiningHandler) CreateDinings(w http.ResponseWriter, r *http.Request){
	
}