package database

import (
	"fmt"
	"github.com/antunesgabriel/babytl_backend/src/infrastructure/models"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(&models.User{}, &models.Album{}, &models.Snap{}, &models.Solicitation{})

	fmt.Println("Run")

	if err != nil {
		fmt.Println(err)
	}
}
