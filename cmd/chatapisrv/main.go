package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"../../pkg/chatapi"
	"github.com/joho/godotenv"
)

var (
	enableWebClient bool             = false
	host            string           = "localhost:8080"
	cAPI            *chatapi.ChatAPI = nil
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	loadEnvVars()

	fmt.Println("Starting server at", host, "...")

	cAPI = chatapi.NewChatAPI()
	go cAPI.Run()

	http.HandleFunc("/", handleRootURL)
	if enableWebClient {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	}
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		chatapi.HandleWebSocket(cAPI, w, r)
	})

	go func() {
		if err := http.ListenAndServe(host, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	fmt.Println("Server started.")

	b, _ := strconv.ParseBool(os.Getenv("ENABLE_CLI"))
	log.Println("ENABLE_CLI", b)
	startCLI(&b)
}

func handleRootURL(writer http.ResponseWriter, request *http.Request) {
	if enableWebClient == false || request.URL.Path != "/" {
		http.Error(writer, "Not found", http.StatusNotFound)
		return
	}
	if request.Method != "GET" {
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(writer, request, "static/home.html")
}

func startCLI(enabled *bool) {
	var input string
	if *enabled {
		reader := bufio.NewReader(os.Stdin)
		for *enabled {
			fmt.Print("ChatAPI > ")

			input, _ = reader.ReadString('\n')

			input = strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r")

			switch input {
			case "msg":
				fmt.Println("test success")
			case "exit", "quit":
				*enabled = false
			default:
				fmt.Println("Command mismatch")
			}

		}
	} else {
		switch {
		}
	}
}

func loadEnvVars() {
	var v, err = strconv.ParseBool(os.Getenv("ENABLE_WEB_CLIENT"))
	if err != nil {
		log.Fatal("ENABLE_WEB_CLIENT parse failed")
	}
	enableWebClient = v
	host = os.Getenv("HOST")
	log.Println("ENABLE_WEB_CLIENT", enableWebClient)
}
