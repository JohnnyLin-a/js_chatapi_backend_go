package main

import (
	"log"

	"../../pkg/chatapi/database"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	database.Migrate()

	log.Println("Execution complete")
}
