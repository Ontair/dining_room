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

type Error struct {
	Error string `json:"error"`
}

// sendError вспомогательная функция для отправки ошибок.
func (d *DishesHandler) sendError(w http.ResponseWriter, msgErr string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := Error{
		Error: msgErr,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		d.logger.Error(
			"Failed to encode error response",
			slog.Any("error", err),
		)
	}
}

func (h *DishesHandler) writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error(
			"Failed to encode error response",
			slog.Any("error", err),
		)
		return err
	}
	return nil
}

func (d *DishesHandler) GetDishes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	dishes, err := d.service.Dishes(ctx)
	if err != nil {
		d.sendError(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := d.writeJSON(w, http.StatusOK, dishes); err != nil {
		d.logger.Error(
			"Response encoding error",
			slog.Any("error", err),
		)
		return
	}
}

func (d *DishesHandler) CreateDish(w http.ResponseWriter, r *http.Request) {
	var data dishData
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		d.logger.Error(
			"Failed to decode JSON",
			slog.Any("error", err),
		)
		d.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx := r.Context()

	dish, err := d.service.CreateDish(ctx, data.Name, data.Price, data.Descriptions)
	if err != nil {
		d.logger.Error(
			"Response encoding error",
			slog.Any("error", err),
		)
		d.sendError(w, err.Error(), http.StatusInternalServerError)

		return
	}

	if err := d.writeJSON(w, http.StatusOK, dish); err != nil {
		d.logger.Error(
			"Response encoding error",
			slog.Any("error", err),
		)
		return
	}
}
