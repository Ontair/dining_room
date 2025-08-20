package http

import (
	"encoding/json"
	"io"
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

func (*DishesHandler) GetDishes(w http.ResponseWriter, r *http.Request) {

}

func (d *DishesHandler) CreateDish(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		d.logger.Error("Не удалось прочитать тело запроса: %v", slog.Any("error", err))
		http.Error(w, "Не удалось прочитать тело запроса", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var data dishData
	err = json.Unmarshal(body, &data)
	if err != nil {
		d.logger.Error("Не удалось декодировать JSON: %v", slog.Any("error", err))
		http.Error(w, "Не удалось декодировать JSON", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	dish, err := d.service.CreateDish(ctx, data.Name, data.Price, data.Descriptions)
	if err != nil {
		d.logger.Error("Ошибка кодирования ответа: %v", slog.Any("error", err))
		http.Error(w, "Ошибка создания блюда", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(dish)
	if err != nil {
		d.logger.Error("Ошибка кодирования ответа: %v", slog.Any("error", err))
	}
}
