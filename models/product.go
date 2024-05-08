package models

import (
	"strings"

	"gorm.io/gorm"
)

const PRODUCTS_PER_PAGE = 10000

type Product struct {
	gorm.Model

	ID          uint
	Slug        string `gorm:"unique"`
	Title       string
	Description string
	Price       int

	CollectionId uint
	Collection   Collection `gorm:"foreignKey:CollectionId"`

	Image      string
	IsFeatured bool
}

func NewProduct(
	slug string,
	title string,
	description string,
	price int,
	collection_id uint,
	image string,
) *Product {
	dashed_slug := strings.ReplaceAll(slug, " ", "-")

	return &Product{
		Slug:         dashed_slug,
		Title:        title,
		Description:  description,
		Price:        price,
		CollectionId: collection_id,
		Image:        image,
		IsFeatured:   false,
	}
}
