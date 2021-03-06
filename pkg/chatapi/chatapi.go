package chatapi

import (
	"encoding/json"
	"log"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database"
	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
)

// ChatAPI manages each client and their actions
type ChatAPI struct {
	messageProcessor chan Message
	register         chan *Client
	unregister       chan *Client
	users            map[uint64]map[*Client]bool
}

// Message is the data structure for incoming messages to the API
type Message struct {
	sender      *Client
	jsonmessage []byte
}

// Response is the data structure for server responses to the client
// Marshall this struct when sending data back to client.
type Response struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Response string `json:"response"`
}

// NewChatAPI creates a new app instance and returns its own pointer
func NewChatAPI() *ChatAPI {
	return &ChatAPI{
		messageProcessor: make(chan Message),
		register:         make(chan *Client),
		unregister:       make(chan *Client),
		users:            make(map[uint64]map[*Client]bool),
	}
}

// Run starts the api main loop to process client data sent from/to server
func (cAPI *ChatAPI) Run() {
	for {
		select {
		case client := <-cAPI.register:
			client.Register()
		case client := <-cAPI.unregister:
			client.Unregister(true)
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

func (cAPI *ChatAPI) broadcastMessage(cMessage *Message) {
	var message models.Message
	var err = json.Unmarshal(cMessage.jsonmessage, &message)
	if err != nil {
		log.Fatalln("BROADCAST: Cannot parse " + string(cMessage.jsonmessage))
		return
	}

	if cMessage.sender.user != nil {
		message.UserID = cMessage.sender.user.ID
		log.Println("#general " + cMessage.sender.user.DisplayName + ": " + message.Message)
	} else {
		message.UserID = 0
		log.Println("#general Guest: " + message.Message)
	}
	

	var db, dbErr = database.NewDatabase()
	if dbErr != nil {
		log.Fatalln("chatapi.broadcastMessage: Database connection failed.")
		return
	}
	message.Save(db)
	db.Close()
	// saveMessage(cMessage.jsonmessage)

	jsonmessage, _ := json.Marshal(message)

	for _, clientsMap := range cAPI.users {
		for client := range clientsMap {
			select {
			case client.send <- jsonmessage:
			default:
				client.Unregister(true)
			}
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
	var chatroom string = ""
	//get last 100 msgs
	var db, err = database.NewDatabase()
	if err != nil {
		log.Fatalln("chatapi.handleOnConnect: Cannot connect to database. ", err)
		return
	}
	jsonMessages, _ := models.GetLast100Messages(db, &chatroom)
	db.Close()

	jsonData, err := json.Marshal(jsonMessages)
	if err != nil {
		log.Println("handleOnConnect: Unable to Marshal jsonMessages")
		cAPI.unregister <- c
		c.conn.Close()
	}
	c.send <- jsonData

	c.SendSysMessage("You are connected to #general.")
}

// Connections gets all connections from the ChatAPi
func (cAPI *ChatAPI) Connections() map[uint64]map[*Client]bool {
	return cAPI.users
}
