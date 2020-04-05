package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	_ "net/http/pprof"

	"../../pkg/chatapi"
	"github.com/joho/godotenv"
)

var (
	enableWebClient bool   = false
	host            string = "localhost:8080"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	loadEnvVars()

	fmt.Println("Starting server at", host, "...")

	cAPI := chatapi.NewChatAPI()
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

	startCLI()
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

func startCLI() {
	// to implement
	select {}
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
