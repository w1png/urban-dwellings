package user_handlers

import (
	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherCollectionsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/collection/:id", CollectionHandler)
	user_api_group.GET("/collection/:id", CollectionApiHandler)
}

func CollectionHandler(c echo.Context) error {
	var collection *models.Collection

	if err := storage.GormStorageInstance.DB.First(&collection, c.Param("id")).Error; err != nil {
		return err
	}

	if err := storage.GormStorageInstance.DB.Where("collection_id = ?", collection.ID).Find(&collection.Products).Error; err != nil {
		return err
	}

	return utils.Render(c, user_templates.Collection(collection))
}

func CollectionApiHandler(c echo.Context) error {
	var collection *models.Collection
	if err := storage.GormStorageInstance.DB.First(&collection, c.Param("id")).Error; err != nil {
		return err
	}
	return utils.Render(c, user_templates.CollectionApi(collection))
}
