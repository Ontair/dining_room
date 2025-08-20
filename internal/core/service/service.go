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

	err := d.repo.Create(ctx, dish)
	return dish, err
}

func (d *DishesService) Dishes(ctx context.Context) ([]*domain.Dish, error) {
	dishes, err := d.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return dishes, err
}
