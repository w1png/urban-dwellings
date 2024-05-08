package user_handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherIndexHandlers(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/", IndexHandler)
	user_api_group.GET("/index", IndexApiHandler)
}

func IndexApiHandler(c echo.Context) error {
	var collections []*models.Collection
	if err := storage.GormStorageInstance.DB.Where("is_enabled = ?", true).Find(&collections).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.IndexApi(collections))
}

func IndexHandler(c echo.Context) error {
	var collections []*models.Collection
	if err := storage.GormStorageInstance.DB.Where("is_enabled = ?", true).Find(&collections).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, user_templates.Index(collections))
}
