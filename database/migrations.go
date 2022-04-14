package database

import (
	"fmt"
	models2 "github.com/antunesgabriel/babytl_backend/src/infrastructure/models"

	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(&models2.User{}, &models2.Album{}, &models2.Snap{}, &models2.Solicitation{})

	fmt.Println("Run")

	if err != nil {
		fmt.Println(err)
	}

}
