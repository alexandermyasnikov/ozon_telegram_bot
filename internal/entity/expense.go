package entity

type Expense struct {
	category string
	price    Decimal
	date     Date
}

func NewExpense(category string, price Decimal, date Date) Expense {
	return Expense{
		category: category,
		price:    price,
		date:     date,
	}
}

func (e *Expense) GetCategory() string {
	return e.category
}

func (e *Expense) GetPrice() Decimal {
	return e.price
}

func (e *Expense) Getdate() Date {
	return e.date
}
