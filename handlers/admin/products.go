package admin_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func GatherProductsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/products", ProductsHandler)
	admin_api_group.GET("/products", ProductsApiHandler)
	admin_api_group.GET("/products/add", AddProductModalHandler)
	admin_api_group.POST("/products", PostProductHandler)
	admin_api_group.GET("/products/:id/edit", EditProductModalHandler)
	admin_api_group.PUT("/products/:id", PutProductHandler)
	admin_api_group.DELETE("/products/:id", DeleteProductHandler)
}

func ProductsHandler(c echo.Context) error {
	search := c.QueryParam("search")

	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Products(products, search))
}

func ProductsApiHandler(c echo.Context) error {
	search := c.QueryParam("search")

	query := storage.GormStorageInstance.DB.Limit(models.PRODUCTS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var products []*models.Product
	if err := query.Find(&products).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.ProductsApi(products, search))
}

func AddProductModalHandler(c echo.Context) error {
	var collections []*models.Collection

	if err := storage.GormStorageInstance.DB.Find(&collections).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.AddProductModal(collections))
}

func PostProductHandler(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(20 << 20); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	price, err := strconv.ParseUint(c.FormValue("price"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	collectionId, err := strconv.ParseUint(c.FormValue("collection"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	image_file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	image, err := image_file.Open()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	processed_image, _, err := utils.ProcessImage(image)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	image_id, err := file_storage.FileStorageInstance.UploadFile(processed_image)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	product := models.NewProduct(
		slug,
		title,
		description,
		int(price),
		uint(collectionId),
		string(image_id),
	)
	if err := storage.GormStorageInstance.DB.Create(product).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Product(product))
}

func EditProductModalHandler(c echo.Context) error {
	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		return err
	}

	var collections []*models.Collection
	if err := storage.GormStorageInstance.DB.Find(&collections).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditProductModal(product, collections))
}

func PutProductHandler(c echo.Context) error {
	var product *models.Product
	if err := storage.GormStorageInstance.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		return err
	}

	if err := c.Request().ParseMultipartForm(20 << 20); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	slug := c.FormValue("slug")
	if slug == "" {
		return c.String(http.StatusBadRequest, "Ссылка не может быть пустой")
	}

	price, err := strconv.ParseUint(c.FormValue("price"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	collectionId, err := strconv.ParseUint(c.FormValue("collection"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	image_file, err := c.FormFile("image")
	var image_id file_storage.ObjectStorageId
	if err == nil {
		image, err := image_file.Open()
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Неверный запрос")
		}

		processed_image, _, err := utils.ProcessImage(image)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Неверный запрос")
		}

		image_id, err = file_storage.FileStorageInstance.UploadFile(processed_image)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Неверный запрос")
		}
	}

	product.Title = title
	product.Slug = slug
	product.Price = int(price)
	product.CollectionId = uint(collectionId)
	product.Description = description
	if image_id != "" {
		product.Image = string(image_id)
	}
	if err := storage.GormStorageInstance.DB.Save(product).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Product(product))
}

func DeleteProductHandler(c echo.Context) error {
	if err := storage.GormStorageInstance.DB.Delete(&models.Product{}, c.Param("id")).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.NoContent(http.StatusOK)
}
