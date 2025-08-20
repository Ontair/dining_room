package service

import (
	"context"
	"log/slog"

	"github.com/Ontair/dining-room/internal/core/ports"
	"github.com/Ontair/dining-room/internal/domain"
	"github.com/brianvoe/gofakeit/v6"
)

var _ ports.DishesService = (*DishesService)(nil)

type DishesService struct {
	repo ports.DishesRepository
	l    *slog.Logger
}

func NewDishesService(repo ports.DishesRepository, l *slog.Logger) *DishesService {
	return &DishesService{
		repo: repo,
		l:    l,
	}
}

func (d *DishesService) CreateDish(ctx context.Context, name, price, descriptions string) (*domain.Dish, error) {
	id := gofakeit.UUID()
	
	dish := domain.NewDishes(id, name, price, descriptions)
	d.l.Debug("создание нового блюда", 
		slog.String("id", id),
		slog.String("name", name),
		slog.String("price", price))

	err := d.repo.Create(ctx, dish)
	if err != nil {
		d.l.Error("ошибка при создании блюда в репозитории", 
			slog.String("error", err.Error()),
			slog.String("id", id))
		return nil, err
	}

	d.l.Debug("блюдо успешно создано", slog.String("id", id))
	return dish, nil
}

func (d *DishesService) Dishes(ctx context.Context) ([]*domain.Dish, error) {
	dishes, err := d.repo.GetAll(ctx)
	if err != nil {
		d.l.Error("ошибка при получении списка блюд из репозитория", 
			slog.String("error", err.Error()))
		return nil, err
	}

	d.l.Info("список блюд успешно получен", 
		slog.Int("count", len(dishes)))
	return dishes, nil
}
