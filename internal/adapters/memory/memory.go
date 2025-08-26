package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/Ontair/dining-room/internal/core/ports"
	"github.com/Ontair/dining-room/internal/domain"
)

var _ ports.DishesRepository = (*MemoryDishesRepository)(nil)

var (
	ErrIsExist = errors.New("ID is already in the map")
)

type MemoryDishesRepository struct {
	dishes map[string]*domain.Dish
	mux    sync.RWMutex
}

func NewMemoryDishesRepository() *MemoryDishesRepository {
	return &MemoryDishesRepository{
		dishes: make(map[string]*domain.Dish),
	}
}

func (d *MemoryDishesRepository) GetAll(_ context.Context) ([]*domain.Dish, error) {
	d.mux.RLock()
	defer d.mux.RUnlock()
	res := make([]*domain.Dish, 0, len(d.dishes))
	for _, dish := range d.dishes {
		res = append(res, dish)
	}
	return res, nil
}

func (d *MemoryDishesRepository) Create(_ context.Context, dish *domain.Dish) error {
	d.mux.Lock()
	defer d.mux.Unlock()

	if i := d.dishes[dish.ID]; i != nil {
		return ErrIsExist
	}
	d.dishes[dish.ID] = dish
	return nil
}
