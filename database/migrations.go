package database

import (
	"fmt"

	"github.com/antunesgabriel/babytl_backend/entities"
	"gorm.io/gorm"
)

func RunMigration(db *gorm.DB) {
	err := db.AutoMigrate(&entities.User{}, &entities.Album{}, &entities.TimeLine{}, &entities.Solicitation{})

	fmt.Println("Run")

	if err != nil {
		fmt.Println(err)
	}

}
