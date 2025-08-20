// ДАННЫЕ (как получить и сохранить).
package ports

import (
	"context"

	"github.com/Ontair/dining-room/internal/domain"
)

type DishesRepository interface {
	GetAll(ctx context.Context) ([]*domain.Dish, error)
	Create(ctx context.Context, dish *domain.Dish) error
}
