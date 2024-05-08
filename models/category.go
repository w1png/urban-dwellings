package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model

	ID   uint
	Name string
	Slug string `gorm:"unique"`
	Tags string

	Products []*Product `gorm:"-"`

	IsEnabled bool
}

const CATEGORIES_PER_PAGE = 100000

func (c *Category) BeforeDelete(tx *gorm.DB) error {
	return tx.Model(&Product{}).Where("category_id = ?", c.ID).Update("category_id", 0).Error
}

func NewCategory(name, slug, tags string) *Category {
	return &Category{
		Name:      name,
		Slug:      slug,
		Tags:      tags,
		IsEnabled: true,
	}
}
