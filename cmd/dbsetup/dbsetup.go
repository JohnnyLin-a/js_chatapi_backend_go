package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	startCLI()
}

func startCLI() {

	var input string
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("dbmanag > ")

		input, _ = reader.ReadString('\n')

		input = strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r")

		switch input {
		case "migrate":
			database.Migrate()
		case "getlast100messages":
			channel := "#general"
			var db, dbErr = database.NewDatabase()
			if dbErr != nil {
				log.Fatalln("dbsetup.getlast100messages: Database connection failed.")
			}
			messages := models.GetLast100Messages(db, &channel)
			for i, message := range messages {
				log.Println(i, message.Timestamp, message.Message)
			}
		case "exit", "quit":
			os.Exit(0)
		default:
			fmt.Println("Command mismatch")
		}

	}
}
