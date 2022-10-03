package entity

import (
	"time"
)

// TODO переименовать product, убрать name.
type Product struct {
	name     string
	category string
	price    int
	date     time.Time
}

func NewProduct(name, category string, price int, date time.Time) Product {
	return Product{
		name:     name,
		category: category,
		price:    price,
		date:     date,
	}
}

func (p *Product) GetName() string {
	return p.name
}

func (p *Product) GetCategory() string {
	return p.category
}

func (p *Product) GetPrice() int {
	return p.price
}

func (p *Product) GetDate() time.Time {
	return p.date
}
