package domain

type Dish struct {
	ID           string
	Name         string
	Price        string
	Descriptions string
}

func NewDishes(id, name, price, descriptions string) *Dish {
	return &Dish{
		ID:           id,
		Name:         name,
		Price:        price,
		Descriptions: descriptions,
	}
}
