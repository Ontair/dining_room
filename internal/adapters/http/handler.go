package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Ontair/dining-room/internal/core/ports"
)

type DishesHandler struct {
	service ports.DishesService
	logger  *slog.Logger
}

type dishData struct {
	Name         string `json:"name"`
	Price        string `json:"price"`
	Descriptions string `json:"description"`
}

func NewDishesHandler(service ports.DishesService, logger *slog.Logger) *DishesHandler {
	return &DishesHandler{
		service: service,
		logger:  logger,
	}
}

// Response представляет стандартный формат ответа.
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

// sendError вспомогательная функция для отправки ошибок.
func (d *DishesHandler) sendError(w http.ResponseWriter, msgErr error, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: false,
		Error:   msgErr,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.logger.Error("Failed to encode error response", "error", err)
	}
}

func (h *DishesHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Response{
		Success: true,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("Failed to encode error response", "error", err)
		return err
	}
	return nil
}

func (d *DishesHandler) GetDishes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dishes, err := d.service.Dishes(ctx)
	if err != nil {
		d.sendError(w, err, http.StatusInternalServerError)
		return
	}
	if err := d.writeJSON(w, http.StatusOK, dishes); err != nil {
		d.logger.Error("Response encoding error", "error", err)
		return
	}
}

func (d *DishesHandler) CreateDish(w http.ResponseWriter, r *http.Request) {
	var data dishData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		d.logger.Error("Failed to decode JSON", "error", err)
		d.sendError(w, err, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	dish, err := d.service.CreateDish(ctx, data.Name, data.Price, data.Descriptions)
	if err != nil {
		d.logger.Error("Response encoding error", "error", err)
		d.sendError(w, err, http.StatusInternalServerError)

		return
	}

	if err := d.writeJSON(w, http.StatusOK, dish); err != nil {
		d.logger.Error("Response encoding error", "error", err)
		return
	}
}
