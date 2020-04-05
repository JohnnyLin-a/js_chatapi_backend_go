package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"../../pkg/chatapi/database"
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
			messages := database.GetLast100Messages("#general")
			for i, message := range messages {
				log.Println(i, message.Timestamp, message.Message)
			}
		case "exit", "quit":
			break
		default:
			fmt.Println("Command mismatch")
		}

	}
}
