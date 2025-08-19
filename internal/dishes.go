package domain

type Dishes struct {
	ID           string
	Name         string
	Price        string
	Descriptions string
}

func NewDishes(id, name, price, descriptions string) *Dishes {
	return &Dishes{
		ID: id,
		Name: name,
		Price: price,
		Descriptions:descriptions,
	}
}
