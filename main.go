package main

import (
	"log"
	"reflect"

	"github.com/w1png/go-htmx-ecommerce-template/config"
	"github.com/w1png/go-htmx-ecommerce-template/errors"
	"github.com/w1png/go-htmx-ecommerce-template/models"
	"github.com/w1png/go-htmx-ecommerce-template/storage"
	"github.com/w1png/go-htmx-ecommerce-template/utils"
)

func createDefaultAdmin() error {
	if err := storage.GormStorageInstance.DB.First(&models.User{}, 1).Error; err == nil {
		return nil
	} else if reflect.TypeOf(err) != reflect.TypeOf(&errors.ObjectNotFoundError{}) {
		return err
	}

	admin, err := models.NewUser("admin", "admin", true)
	if err != nil {
		return err
	}

	return storage.GormStorageInstance.DB.Create(&admin).Error
}

func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}

	if err := storage.InitStorage(); err != nil {
		log.Fatal(err)
	}

	if err := createDefaultAdmin(); err != nil {
		log.Fatal(err)
	}

	if err := utils.CreateImagesFolderIfNotExists(); err != nil {
		log.Fatal(err)
	}

	server := NewHTTPServer()

	log.Fatal(server.Run())
}
