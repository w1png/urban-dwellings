package models

import "gorm.io/gorm"

type OrderProduct struct {
	gorm.Model

	ID        uint
	OrderID   uint
	ProductId uint
	Slug      string
	Name      string
	Price     int
	Quantity  int
}

func NewOrderProduct(
	product_id uint,
	order_id uint,
	slug string,
	name string,
	price int,
	quantity int,
) *OrderProduct {
	return &OrderProduct{
		ProductId: product_id,
		OrderID:   order_id,
		Slug:      slug,
		Name:      name,
		Price:     price,
		Quantity:  quantity,
	}
}

func (op *OrderProduct) GetTotalPrice() int {
	p := op.Price
	return p * op.Quantity
}
