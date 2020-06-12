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
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: messageSplit[0], Response: cMessage.sender.displayName}
		var jsonResponse, err = json.Marshal(jsonResponseStruct)
		if err != nil {
			log.Println("Marshal error", jsonData["message"], cMessage.sender.displayName, err)
			return
		}
		cMessage.sender.send <- jsonResponse
		return
	case "!register":
		register(&messageSplit, cMessage.sender)
	default:
		log.Println("Unable to process SYSCOMMAND:", jsonData["message"])
		return
	}

}

func register(message *[]string, sender *Client) {
	// !register email username displayName password
	if len(*message) != 5 {
		log.Fatalln("sys_commands.register: failed. args:", *message)
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
