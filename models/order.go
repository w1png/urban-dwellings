package models

import (
	"time"

	"gorm.io/gorm"
)

const ORDERS_PER_PAGE = 200000

type Order struct {
	gorm.Model

	ID          uint
	Name        string
	PhoneNumber string
	Email       string
	Message     string
	Products    []*OrderProduct
	IsResolved  bool
}

func NewOrder(
	name string,
	phone_number string,
	email string,
	message string,
) *Order {
	return &Order{
		Name:        name,
		PhoneNumber: phone_number,
		Email:       email,
		Message:     message,
		IsResolved:  false,
	}
}

func (o *Order) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(&OrderProduct{}).Where("order_id = ?", o.ID).Find(&o.Products).Error; err != nil {
		return err
	}
	return nil
}

func (o *Order) GetTotalPrice() int {
	total := 0
	for _, product := range o.Products {
		p := product.Price
		total += p * product.Quantity
	}

	return total
}

func (o *Order) FormatTime() string {
	return o.CreatedAt.In(time.Local).Format("15:04 02-01-2006")
}
