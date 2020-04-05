package chatapi

import (
	"encoding/json"
	"log"
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
	sender  *Client
	message []byte
}

// Response is the data structure for server responses to the client
// Marshall this struct when sending data back to client.
type Response struct {
	Type     string
	Message  string
	Response string
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
				cAPI.broadcastMessage(&Message{client, []byte(`{"Type":"MESSAGE","Message":"` + uname + ` disconnected.","Sender":"SYSTEM"}`)})
			}
		case cMessage := <-cAPI.messageProcessor:
			go cAPI.processMessage(cMessage)
		}
	}
}

func (cAPI *ChatAPI) processMessage(cMessage Message) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(cMessage.message, &jsonData); err != nil {
		log.Println("UNMARSHAL ERROR ", err)
		log.Println("for json: ", cMessage.message)
		return
	}

	switch jsonData["Type"] {
	case "_SYSCOMMAND":
		cAPI.processSysCommand(&cMessage)
	case "MESSAGE":
		cAPI.broadcastMessage(&cMessage)
	default:
		log.Println("Unable to process message of type", jsonData["Type"])
		return
	}
}

func (cAPI *ChatAPI) processSysCommand(cMessage *Message) {
	log.Println("Process SYSCOMMAND", string(cMessage.message))

	jsonData, err := parseJSON(cMessage.message)
	if err != nil {
		return
	}

	switch jsonData["Message"] {
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
	var parsedMessage, err = parseJSON(cMessage.message)
	if err != nil {
		log.Fatalln("BROADCAST: Cannot parse " + string(cMessage.message))
		return
	}
	log.Println("#general " + parsedMessage["Sender"].(string) + ": " + parsedMessage["Message"].(string))
	saveMessage(cMessage.message)
	for client := range cAPI.clients {
		select {
		case client.send <- cMessage.message:
		default:
			uname := client.displayName
			close(client.send)
			delete(cAPI.clients, client)
			cAPI.broadcastMessage(&Message{client, []byte(`{"type":"MESSAGE","message":"` + uname + ` disconnected.","Sender":"SYSTEM"}`)})

		}
	}
}

func parseJSON(message []byte) (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	if err := json.Unmarshal(message, &jsonData); err != nil {
		log.Println("UNMARSHAL ERROR ", err)
		log.Println("for json: ", string(message))
		return nil, err
	}
	return jsonData, nil
}
