package chatapi

import (
	"encoding/json"
	"log"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
)

// ChatAPI manages each client and their actions
type ChatAPI struct {
	clients          map[*Client]bool
	messageProcessor chan Message
	register         chan *Client
	unregister       chan *Client
}

// Message is the data structure for incoming messages to the API
type Message struct {
	sender      *Client
	jsonmessage []byte
}


// Response is the data structure for server responses to the client
// Marshall this struct when sending data back to client.
type Response struct {
	Type     string	`json:"type"`
	Message  string	`json:"message"`
	Response string	`json:"response"`
}

// NewChatAPI creates a new app instance and returns its own pointer
func NewChatAPI() *ChatAPI {
	return &ChatAPI{
		messageProcessor: make(chan Message),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		clients:          make(map[*Client]bool),
	}
}

// Run starts the api main loop to process client data sent from/to server
func (cAPI *ChatAPI) Run() {
	for {
		select {
		case client := <-cAPI.register:
			cAPI.clients[client] = true
		case client := <-cAPI.unregister:
			if _, ok := cAPI.clients[client]; ok {
				uname := client.displayName
				close(client.send)
				delete(cAPI.clients, client)
				cAPI.broadcastMessage(&Message{client, []byte(`{"type":"MESSAGE","message":"` + uname + ` disconnected.","sender":"SYSTEM"}`)})
			}
		case cMessage := <-cAPI.messageProcessor:
			go cAPI.processMessage(cMessage)
		}
	}
}

func (cAPI *ChatAPI) processMessage(cMessage Message) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(cMessage.jsonmessage, &jsonData); err != nil {
		log.Println("UNMARSHAL ERROR ", err)
		log.Println("for json: ", cMessage.jsonmessage)
		return
	}

	switch jsonData["type"] {
	case "_SYSCOMMAND":
		cAPI.processSysCommand(&cMessage)
	case "MESSAGE":
		cAPI.broadcastMessage(&cMessage)
	default:
		log.Println("Unable to process message of type", jsonData["type"])
		return
	}
}

func (cAPI *ChatAPI) processSysCommand(cMessage *Message) {
	log.Println("Process SYSCOMMAND", string(cMessage.jsonmessage))

	jsonData, err := parseGenericJSON(cMessage.jsonmessage)
	if err != nil {
		return
	}

	switch jsonData["message"] {
	case "!get_display_name":
		var jsonResponseStruct = Response{Type: "_SYSCOMMAND", Message: "!get_display_name", Response: cMessage.sender.displayName}
		var jsonResponse, err = json.Marshal(jsonResponseStruct)
		if err != nil {
			log.Println("Marshal error", jsonData["message"], cMessage.sender.displayName, err)
			return
		}
		cMessage.sender.send <- jsonResponse

	default:
		log.Println("Unable to process SYSCOMMAND:", jsonData["message"])
		return
	}

}

func (cAPI *ChatAPI) broadcastMessage(cMessage *Message) {
	var message models.Message
	var err = json.Unmarshal(cMessage.jsonmessage, &message)
	if err != nil {
		log.Fatalln("BROADCAST: Cannot parse " + string(cMessage.jsonmessage))
		return
	}
	log.Println("#general " + message.Sender + ": " + message.Message)

	var db, dbErr = database.NewDatabase()
	if dbErr != nil {
		log.Fatalln("chatapi.broadcastMessage: Database connection failed.")
		return
	}
	message.SaveMessage(db)
	// db connection will be closed here
	// saveMessage(cMessage.jsonmessage)

	for client := range cAPI.clients {
		select {
		case client.send <- cMessage.jsonmessage:
		default:
			uname := client.displayName
			close(client.send)
			delete(cAPI.clients, client)
			cAPI.broadcastMessage(&Message{client, []byte(`{"type":"MESSAGE","message":"` + uname + ` disconnected.","sender":"SYSTEM"}`)})

		}
	}
}

func parseGenericJSON(message []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(message, &jsonData); err != nil {
		log.Println("UNMARSHAL ERROR ", err)
		log.Println("for json: ", string(message))
		return nil, err
	}
	return jsonData, nil
}

func (cAPI *ChatAPI) handleOnConnect(c *Client) {
	var chatroom string = "#general"
	//get last 100 msgs
	var db, err = database.NewDatabase()
	if err != nil {
		log.Fatalln("chatapi.handleOnConnect: Cannot connect to database. ", err)
		return
	}
	jsonMessages := models.GetLast100Messages(db, &chatroom)
	// db connection will close here.

	jsonData, err := json.Marshal(jsonMessages)
	if err != nil {
		log.Println("handleOnConnect: Unable to Marshal jsonMessages")
		cAPI.unregister <- c
		c.conn.Close()
	}
	c.send <- jsonData

	// broadcast new client connection
	cAPI.broadcastMessage(&Message{c, []byte(`{"type":"MESSAGE","message":"` + c.displayName + ` connected to #general.","sender":"SYSTEM"}`)})
}
