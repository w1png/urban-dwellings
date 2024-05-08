package admin_handlers

import (
	"github.com/labstack/echo"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherOrderRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/orders", OrdersHandler)
	admin_api_group.GET("/orders", OrdersApiHandler)
}

func OrdersHandler(c echo.Context) error {
	var orders []*models.Order
	if err := storage.GormStorageInstance.DB.Find(&orders).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.Orders(orders))
}

func OrdersApiHandler(c echo.Context) error {
	var orders []*models.Order
	if err := storage.GormStorageInstance.DB.Find(&orders).Error; err != nil {
		return err
	}

	return utils.Render(c, admin_templates.OrdersApi(orders))
}
