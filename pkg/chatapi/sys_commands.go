package chatapi

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
)

func (cAPI *ChatAPI) processSysCommand(cMessage *Message) {

	jsonData, err := parseGenericJSON(cMessage.jsonmessage)
	if err != nil {
		return
	}

	messageSplit := strings.Split((jsonData["message"].(string)), " ")
	log.Println("PROCESSING: ", messageSplit[0])

	switch messageSplit[0] {
	case "!get_display_name":
		getDisplayName(&messageSplit, cMessage.sender)
	case "!register":
		register(&messageSplit, cMessage.sender)
	case "!login":
		login(&messageSplit, cMessage.sender)
	default:
		log.Println("Unable to process SYSCOMMAND:", jsonData["message"])
		return
	}

}

func getDisplayName(message *[]string, sender *Client) {
	var displayName string
	if sender.user.ID == 0 {
		displayName = "Guest"
	} else {
		displayName = (*sender).user.DisplayName
	}
	var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: (*message)[0], Response: displayName}
	var jsonResponse, _ = json.Marshal(jsonResponseStruct)

	sender.send <- jsonResponse
	return
}

func register(message *[]string, sender *Client) {
	// !register email username displayName password
	if len(*message) != 5 {
		log.Println("sys_commands.register: failed. args count:", len(*message))
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!register", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}

	var db, _ = database.NewDatabase()
	var hashedPassword, _ = models.Hash((*message)[4])
	var u = models.User{CreatedAt: time.Now(), UpdatedAt: time.Now(), Email: (*message)[1], Username: (*message)[2], DisplayName: (*message)[3], Password: string(hashedPassword)}
	var _, err = u.SaveUser(db)
	if err != nil {
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!register", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}
	db.Close()

	var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!register", Response: "SUCCESS"}
	var jsonResponse, _ = json.Marshal(jsonResponseStruct)

	log.Println("!register SUCCESS")

	sender.send <- jsonResponse
}

func login(message *[]string, sender *Client) {
	if len(*message) != 3 {
		log.Println("sys_commands.login: failed. args:", *message)
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}

	// TODO: login logic
}
