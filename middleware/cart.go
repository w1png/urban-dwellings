package middleware

import (
	"context"
	"net/http"

	"gorm.io/gorm"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
)

func UseCart(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cart_uuid, err := c.Cookie("cart_uuid")
		if err != nil && err != http.ErrNoCookie {
			log.Error(err)
			return err
		}
		var cart *models.Cart

		if cart_uuid == nil {
			cart = models.NewCart()
			if err := storage.GormStorageInstance.DB.Create(&cart).Error; err != nil {
				return err
			}

			c.SetCookie(&http.Cookie{
				Name:  "cart_uuid",
				Path:  "/",
				Value: cart.UUID.String(),
			})
		} else {
			uuid, err := uuid.Parse(cart_uuid.Value)
			if err != nil {
				return err
			}
			if err := storage.GormStorageInstance.DB.Where("uuid = ?", uuid).First(&cart).Error; err != nil && err != gorm.ErrRecordNotFound {
				return err
			}
		}

		c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), "cart", cart)))
		return next(c)
	}
}
