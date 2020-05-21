package chatapi

import (
	"log"
	"time"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
)

func saveMessage(message []byte) {
	var jsonMap, err = parseJSON(message)
	if err != nil {
		log.Println("dbtransaction.saveMessage: Error parsing json")
	}

	var m = database.Message{Chatroom: "#general", Timestamp: time.Now(), Sender: jsonMap["Sender"].(string), Type: jsonMap["Type"].(string), Message: jsonMap["Message"].(string)}
	database.SaveMessage(&m)
}

func getlast100Messages(chatroom *string) *[]JSONMessage {
	dbmessages := database.GetLast100Messages(chatroom)

	var jsonMessages []JSONMessage

	for _, dbmessage := range dbmessages {
		jsonMessages = append(jsonMessages, JSONMessage{Type: dbmessage.Type, Sender: dbmessage.Sender, Message: dbmessage.Message, Timestamp: dbmessage.Timestamp.String(), Chatroom: dbmessage.Chatroom})
	}
	return &jsonMessages
}
