package models

import (
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"gorm.io/gorm"
)

const COLLECTIONS_PER_PAGE = 100000

type Collection struct {
	gorm.Model

	ID          uint
	Title       string
	Description string

	Image     file_storage.ObjectStorageId
	Thumbnail file_storage.ObjectStorageId

	Products []*Product `gorm:"-"`

	IsEnabled bool
}

func NewCollection(title, description string, image, thumbnail file_storage.ObjectStorageId, is_enabled bool) *Collection {
	return &Collection{
		Title:       title,
		Description: description,
		Image:       image,
		Thumbnail:   thumbnail,
		IsEnabled:   is_enabled,
	}
}
