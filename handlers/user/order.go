package user_handlers

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherOrdersRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_api_group.POST("/checkout", PostOrderHandler)
}

func PostOrderHandler(c echo.Context) error {
	c.Response().Header().Set("HX-Reswap", "innerHTML")

	if err := c.Request().ParseForm(); err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	name := c.FormValue("name")
	if name == "" {
		return c.String(http.StatusBadRequest, "Имя не может быть пустым")
	}

	message := c.FormValue("name")
	phone_number := c.FormValue("phone_number")
	if !utils.ValidatePhoneNumber(phone_number) {
		return c.String(http.StatusBadRequest, "Неправильный формат номера телефона")
	}
	email := c.FormValue("email")
	if !utils.ValidateEmail(email) {
		return c.String(http.StatusBadRequest, "Неправильный формат адреса электронной почты")
	}

	cart := utils.GetCartFromContext(c.Request().Context())
	if len(cart.Products) == 0 {
		return c.String(http.StatusBadRequest, "Корзина пуста")
	}

	order := models.NewOrder(
		name,
		phone_number,
		email,
		message,
	)

	if err := storage.GormStorageInstance.DB.Create(order).Error; err != nil {
		log.Error(err)
		return err
	}

	var order_products []*models.OrderProduct
	for _, cart_product := range cart.Products {
		if cart_product.Quantity == 0 {
			continue
		}

		order_product := models.NewOrderProduct(
			cart_product.Product.ID,
			order.ID,
			cart_product.Slug,
			cart_product.Title,
			cart_product.Price,
			cart_product.Quantity,
		)

		if err := storage.GormStorageInstance.DB.Create(order_product).Error; err != nil {
			log.Error(err)
			return err
		}

		cart_product.Quantity = 0
		if err := storage.GormStorageInstance.DB.Save(cart_product).Error; err != nil {
			log.Error(err)
			return err
		}

		order_products = append(order_products, order_product)
	}

	order.Products = order_products

	c.Response().Header().Del("HX-Reswap")

	return utils.Render(c, user_templates.CheckoutComplete())
}
