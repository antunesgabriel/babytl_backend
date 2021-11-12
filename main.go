package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/antunesgabriel/babytl_backend/api"
	"github.com/antunesgabriel/babytl_backend/database"
)

func main() {
	isHeroku := os.Getenv("BTL_SETUP")

	if isHeroku == "" {
		err := godotenv.Load()

		if err != nil {
			log.Fatalln(err)
		}
	}

	database.StartDabate()

	server := api.BuildServer()

	server.Run()
}
