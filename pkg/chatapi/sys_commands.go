package chatapi

import (
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
	"github.com/badoux/checkmail"
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
	case "!logout":
		logout(cMessage.sender)
	default:
		log.Println("Unable to process SYSCOMMAND:", jsonData["message"])
		return
	}

}

func getDisplayName(message *[]string, sender *Client) {
	var displayName string
	if sender.user == nil {
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
	// !login email/username password
	if len(*message) != 3 {
		log.Println("sys_commands.login: failed. args count:", len(*message))
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}

	isEmail := true
	if err := checkmail.ValidateFormat((*message)[1]); err != nil {
		isEmail = false
	}

	var db, _ = database.NewDatabase()
	u := models.User{}
	if isEmail {
		db.First(&u, "email = ?", (*message)[1])
		log.Println("Login with email", (*message)[1])
	} else {
		db.First(&u, "username = ?", (*message)[1])
		log.Println("Login with username", (*message)[1])
	}
	db.Close()

	if u.ID == 0 {
		log.Println("sys_commands.login: failed. login DNE.")
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}

	validPassword := true
	if err := models.VerifyPassword(u.Password, (*message)[2]); err != nil {
		validPassword = false
	}
	log.Println("Password valid", validPassword)

	if !validPassword {
		log.Println("sys_commands.login: failed. Password mismatch")
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}
	sender.user = &u
	log.Println("Login: success")
	var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "SUCCESS"}
	var jsonResponse, _ = json.Marshal(jsonResponseStruct)

	sender.send <- jsonResponse

	jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!get_display_name", Response: u.DisplayName}
	jsonResponse, _ = json.Marshal(jsonResponseStruct)

	sender.send <- jsonResponse
	return
}

func logout(sender *Client) {
	if sender.user == nil {
		log.Println("sys_commands.logout: failed. Not logged in")
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "FAILED"}
		var jsonResponse, _ = json.Marshal(jsonResponseStruct)

		sender.send <- jsonResponse
		return
	}
	log.Println("Logout: ", sender.user.DisplayName)
	sender.user = nil
	var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!login", Response: "SUCCESS"}
	var jsonResponse, _ = json.Marshal(jsonResponseStruct)

	sender.send <- jsonResponse
	return

}
