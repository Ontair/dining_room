package memory

import (
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
	dishes map[string]*domain.Dish // map[string]*domain.
	mux    sync.RWMutex
}

func NewMemoryDishesRepository() *MemoryDishesRepository {
	return &MemoryDishesRepository{
		dishes: make(map[string]*domain.Dish),
	}
}

func (d *MemoryDishesRepository) GetAll() ([]*domain.Dish, error) {
	d.mux.Lock()
	defer d.mux.Unlock()
	res := make([]*domain.Dish, len(d.dishes))
	for _, dish := range d.dishes {
		res = append(res, dish)
	}
	return res, nil
}

func (d *MemoryDishesRepository) Create(dish *domain.Dish) error{
	d.mux.Lock()
	defer d.mux.Unlock()

	if i := d.dishes[dish.ID]; i!=nil{
		return ErrIsExist
	}
	d.dishes[dish.ID] = dish
	return nil
}
