package migrations

import (
	"go-boilerplate/db"
	"go-boilerplate/models"
	"log"
)

func Migrate() {

	err := db.GetDB().AutoMigrate(&models.User{}, &models.Campaign{}, &models.Subscribers{})

	if err != nil {
		log.Panic(err)
	}
}
