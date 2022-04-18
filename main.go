package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	isHeroku := os.Getenv("BTL_SETUP")

	if isHeroku == "" {
		err := godotenv.Load()

		if err != nil {
			log.Fatalln(err)
		}
	}

}
