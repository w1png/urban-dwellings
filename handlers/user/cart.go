package user_handlers

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/templates/components"
	user_templates "github.com/w1png/go-htmx-ecommerce-template/templates/user"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherCartRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	user_page_group.GET("/cart", GetCartHandler)
	user_api_group.GET("/cart", GetCartApiHandler)
	user_api_group.PUT("/cart/change_quantity/:product_id", ChangeCartProductQuantityHandler)
}

func GetCartHandler(c echo.Context) error {
	var cart_products []*models.CartProduct
	for _, cart_product := range utils.GetCartFromContext(c.Request().Context()).Products {
		if cart_product.Quantity != 0 {
			cart_products = append(cart_products, cart_product)
		}
	}

	return utils.Render(c, user_templates.Cart(cart_products))
}

func GetCartApiHandler(c echo.Context) error {
	var cart_products []*models.CartProduct
	for _, cart_product := range utils.GetCartFromContext(c.Request().Context()).Products {
		if cart_product.Quantity != 0 {
			cart_products = append(cart_products, cart_product)
		}
	}

	return utils.Render(c, user_templates.CartApi(cart_products))
}

func ChangeCartProductQuantityHandler(c echo.Context) error {
	should_decrease := c.QueryParam("decrease") == "true"

	product_id, err := strconv.ParseUint(c.Param("product_id"), 10, 64)
	if err != nil {
		return err
	}

	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("id = ?", product_id).First(&product).Error; err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Товар не найден")
		}
		log.Error(err)
		return err
	}

	cart := utils.GetCartFromContext(c.Request().Context())

	var cart_product *models.CartProduct
	if err := storage.GormStorageInstance.DB.Where("product_id = ? AND cart_id = ?", product_id, cart.ID).First(&cart_product).Error; err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&errors.ObjectNotFoundError{}) {
			return c.String(http.StatusNotFound, "Товар не найден")
		}

		cart_product = models.NewCartProduct(
			product.ID,
			cart.ID,
			product.Slug,
			product.Title,
			product.Price,
			0,
		)
	}

	if should_decrease && cart_product.Quantity > 0 {
		cart_product.Quantity--
	} else {
		cart_product.Quantity++
	}

	if err := storage.GormStorageInstance.DB.Save(&cart_product).Error; err != nil {
		log.Error(err)
		return err
	}

	return utils.Render(c, components.AddToCartButton(cart_product.Product.ID, cart_product.Quantity))
}
