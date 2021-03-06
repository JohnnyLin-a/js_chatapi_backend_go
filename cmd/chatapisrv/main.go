package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi"
	"github.com/joho/godotenv"
)

var (
	enableWebClient bool   = false
	httpHost        string = ":8080"
	// httpsHost       string           = ":443"
	dnsAddrOrigin string           = ""
	cAPI          *chatapi.ChatAPI = nil
	// sslCert         string           = "fullchain.crt"
	// sslKey          string           = ""
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	loadEnvVars()

	fmt.Println("Starting server at", httpHost)

	cAPI = chatapi.NewChatAPI()
	go cAPI.Run()

	// http to https redirector
	// http.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir("/.well-known/acme-challenge/"))))
	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	// 	http.Redirect(w, r, "https://"+r.Host, http.StatusPermanentRedirect)
	// })
	// go http.ListenAndServe(httpHost, nil)
	// now using nginx to do that

	// Main app portion
	// mux := http.NewServeMux()

	// mux config
	// cfg := &tls.Config{
	// 	MinVersion:               tls.VersionTLS12,
	// 	CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	// 	PreferServerCipherSuites: true,
	// 	CipherSuites: []uint16{
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
	// 		tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
	// 		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
	// 	},
	// }
	// srv := &http.Server{
	// 	Addr:         httpsHost,
	// 	Handler:      mux,
	// 	TLSConfig:    cfg,
	// 	TLSNextProto: make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0),
	// }

	http.HandleFunc("/", handleRootURL)
	if enableWebClient {
		http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	}
	http.Handle("/.well-known/acme-challenge/", http.StripPrefix("/.well-known/acme-challenge/", http.FileServer(http.Dir("/.well-known/acme-challenge/"))))
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		log.Println("User connecting from " + GetIP(r))
		if r.Header["Origin"][0] != dnsAddrOrigin {
			log.Println("CSRF mismatch origin: " + r.Header["Origin"][0])
		} else {
			chatapi.HandleWebSocket(cAPI, w, r)
		}
	})

	go func() {
		// if _, err := os.Stat(sslCert); err == nil {
		if err := http.ListenAndServe(httpHost, nil); err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
		// } else if os.IsNotExist(err) {
		// log.Println("SSL cert not found. Will not listen on " + httpsHost)
		// }
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
			case "test":
				fmt.Println("test success")
			case "clients":
				fmt.Println(cAPI.Connections())
			case "exit", "quit":
				*enabled = false
			default:
				fmt.Println("Command mismatch")
			}
		}
	} else {
		for {
		}
	}
}

func loadEnvVars() {
	var v, err = strconv.ParseBool(os.Getenv("ENABLE_WEB_CLIENT"))
	if err != nil {
		log.Fatal("ENABLE_WEB_CLIENT parse failed")
		os.Exit(1)
	}
	enableWebClient = v
	httpHost = os.Getenv("HTTP_HOST")
	// httpsHost = os.Getenv("HTTPS_HOST")
	// sslCert = os.Getenv("SSL_CERT")
	// sslKey = os.Getenv("SSL_KEY")
	dnsAddrOrigin = os.Getenv("DNS_ADDR_ORIGIN")
	log.Println("ENABLE_WEB_CLIENT", enableWebClient)
}

// GetIP returns the http request's remote IP address
func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}
