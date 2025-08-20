// ДАННЫЕ (как получить и сохранить).
package ports

import "github.com/Ontair/dining-room/internal/domain"


type DishesRepository interface{
	GetAll() ([]*domain.Dish, error)
	Create(dish *domain.Dish) error
}