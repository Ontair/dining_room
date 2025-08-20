// БИЗНЕС-ЛОГИКУ (что сделать с этими данными)
package ports

import (
	"context"

	"github.com/Ontair/dining-room/internal/domain"
)


type DishesService interface{
	CreateDish(ctx context.Context ,name, price, descriptions string) (*domain.Dish, error)
	Dishes()
}