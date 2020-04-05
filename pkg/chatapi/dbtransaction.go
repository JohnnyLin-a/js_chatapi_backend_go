package chatapi

import (
	"log"
	"time"

	"./database"
)

func saveMessage(message []byte) {
	var jsonMap, err = parseJSON(message)
	if err != nil {
		log.Println("dbtransaction.saveMessage: Error parsing json")
	}

	var m = database.Message{Chatroom: "#general", Timestamp: time.Now(), Sender: jsonMap["Sender"].(string), Type: jsonMap["Type"].(string), Message: jsonMap["Message"].(string)}
	database.SaveMessage(&m)
}
