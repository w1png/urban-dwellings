package admin_handlers

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/w1png/go-htmx-ecommerce-template/file_storage"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/utils"

	"github.com/w1png/go-htmx-ecommerce-template/models"
	admin_templates "github.com/w1png/go-htmx-ecommerce-template/templates/admin"
)

func GatherCollectionsRoutes(user_page_group *echo.Echo, user_api_group, admin_page_group, admin_api_group *echo.Group) {
	admin_page_group.GET("/collections", CollectionsHandler)
	admin_api_group.GET("/collections", CollectionsApiHandler)
	admin_api_group.GET("/collections/page/:page", CollectionsPageHandler)
	admin_api_group.POST("/collections/search", CollectionsSearchHandler)
	admin_api_group.GET("/collections/add", AddCollectionModalHandler)
	admin_api_group.POST("/collections", PostCollectionHandler)
	admin_api_group.GET("/collections/:id/delete", DeleteCollectionModalHandler)
	admin_api_group.DELETE("/collections/:id", DeleteCollectionHandler)
	admin_api_group.GET("/collections/:id/edit", EditCollectionModalHandler)
	admin_api_group.PUT("/collections/:id", PutCollectionHandler)
}

func CollectionsHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.COLLECTIONS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}
	var collections []*models.Collection
	if err := query.Find(&collections).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	return utils.Render(c, admin_templates.Collections(collections, search))
}

func CollectionsApiHandler(c echo.Context) error {
	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.COLLECTIONS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}
	var collections []*models.Collection
	if err := query.Find(&collections).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	return utils.Render(c, admin_templates.CollectionsApi(collections, search))
}

func CollectionsPageHandler(c echo.Context) error {
	page, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	search := c.QueryParam("search")
	query := storage.GormStorageInstance.DB.Limit(models.COLLECTIONS_PER_PAGE).Offset((page - 1) * models.COLLECTIONS_PER_PAGE)
	if search != "" {
		query = query.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%")
	}

	var collections []*models.Collection
	if err := query.Find(&collections).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}
	return utils.Render(c, admin_templates.CollectionsList(collections, page+1, search))
}

func CollectionsSearchHandler(c echo.Context) error {
	if err := c.Request().ParseForm(); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	search := c.FormValue("search")
	var collections []*models.Collection
	if search == "" {
		c.Response().Header().Set("HX-Replace-Url", "/admin/collections")
		if err := storage.GormStorageInstance.DB.Find(&collections).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	} else {
		c.Response().Header().Set("HX-Replace-Url", "/admin/collections?search="+search)
		if err := storage.GormStorageInstance.DB.Where("LOWER(title) LIKE LOWER(?)", "%"+search+"%").Find(&collections).Error; err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	return utils.Render(c, admin_templates.CollectionsList(collections, 2, search))
}

func AddCollectionModalHandler(c echo.Context) error {
	return utils.Render(c, admin_templates.AddCollectionModal())
}

func PostCollectionHandler(c echo.Context) error {
	if err := c.Request().ParseMultipartForm(20 << 20); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	image_file, err := c.FormFile("image")
	if err != nil {
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	image, err := image_file.Open()
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	processed_image, thumbnail, err := utils.ProcessImage(image)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	image_id, err := file_storage.FileStorageInstance.UploadFile(processed_image)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	thumbnail_id, err := file_storage.FileStorageInstance.UploadFile(thumbnail)
	if err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	collection := models.NewCollection(
		title,
		description,
		file_storage.ObjectStorageId(image_id),
		file_storage.ObjectStorageId(thumbnail_id),
		is_enabled,
	)

	if err := storage.GormStorageInstance.DB.Create(&collection).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Collection(collection))
}

func DeleteCollectionModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var collection *models.Collection
	if err := storage.GormStorageInstance.DB.First(&collection, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.DeleteCollectionModal(collection))
}

func DeleteCollectionHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := storage.GormStorageInstance.DB.Delete(&models.Collection{}, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return c.NoContent(http.StatusOK)
}

func EditCollectionModalHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	var collection *models.Collection
	if err := storage.GormStorageInstance.DB.First(&collection, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.EditCollectionModal(collection))
}

func PutCollectionHandler(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Неверный запрос")
	}

	if err := c.Request().ParseMultipartForm(20 << 20); err != nil {
		log.Error(err)
		return c.String(http.StatusBadRequest, "Неправильный запрос")
	}

	title := c.FormValue("title")
	if title == "" {
		return c.String(http.StatusBadRequest, "Название не может быть пустым")
	}

	description := c.FormValue("description")
	if description == "" {
		return c.String(http.StatusBadRequest, "Описание не может быть пустым")
	}

	is_enabled, _ := strconv.ParseBool(c.FormValue("is_enabled"))

	var image_id, thumbnail_id file_storage.ObjectStorageId
	image_file, err := c.FormFile("image")
	if err == nil {
		image, err := image_file.Open()
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Неправильный запрос")
		}

		processed_image, thumbnail, err := utils.ProcessImage(image)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusBadRequest, "Неправильный запрос")
		}

		image_id, err = file_storage.FileStorageInstance.UploadFile(processed_image)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}

		thumbnail_id, err = file_storage.FileStorageInstance.UploadFile(thumbnail)
		if err != nil {
			log.Error(err)
			return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
		}
	}

	var collection *models.Collection
	if err := storage.GormStorageInstance.DB.First(&collection, id).Error; err != nil {
		log.Error(err)
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	collection.Title = title
	collection.Description = description
	if image_id != "" {
		collection.Image = file_storage.ObjectStorageId(image_id)
		collection.Thumbnail = file_storage.ObjectStorageId(thumbnail_id)
	}
	collection.IsEnabled = is_enabled

	if err := storage.GormStorageInstance.DB.Save(&collection).Error; err != nil {
		return c.String(http.StatusInternalServerError, "Неизвестная ошибка")
	}

	return utils.Render(c, admin_templates.Collection(collection))
}
