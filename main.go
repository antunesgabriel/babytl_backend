package main

import (
	"github.com/antunesgabriel/babytl_backend/api"
	"github.com/antunesgabriel/babytl_backend/database"
)

func main() {
	database.StartDabate()

	server := api.BuildServer()

	server.Run()
}
