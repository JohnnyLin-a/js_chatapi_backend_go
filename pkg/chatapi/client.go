package chatapi

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
	"github.com/gorilla/websocket"
)

const (
	writeTimeout = 10 * time.Second
	readTimeout  = 60 * time.Second
	pingInterval = (readTimeout * 9) / 10

	// Maximum message length allowed.
	// length 2000 + 1000 JSON headroom
	maxMessageSize = 3000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// var (
// 	userNumber int64 = 1
// )

// Client is a bridge between the websocket connection and the ChatAPI.
type Client struct {
	cAPI *ChatAPI
	conn *websocket.Conn
	send chan []byte
	user *models.User
}

// startWebsocketReader reads socket's incoming messages to the server
func (c *Client) startWebsocketReader() {
	// Read user websocket input incoming to server

	defer func() {
		c.cAPI.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(readTimeout))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("UnexpectedCloseError:", err)
			}
			break
		}

		// message manipulation. To edit in later versions.
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		// message = append([]byte(c.displayName+": "), message...)

		c.cAPI.messageProcessor <- Message{sender: c, jsonmessage: message}
	}
}

// startWebsocketWriter sends messages from the chatAPI to the connected websocket.
func (c *Client) startWebsocketWriter() {
	// Write to client's socket
	interval := time.NewTicker(pingInterval)

	defer func() {
		interval.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-interval.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeTimeout))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// HandleWebSocket handles client websocket
func HandleWebSocket(cAPI *ChatAPI, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{cAPI: cAPI, conn: conn, send: make(chan []byte, 256)}

	client.cAPI.register <- client

	go client.startWebsocketReader()
	go client.startWebsocketWriter()

	go cAPI.handleOnConnect(client)
}

// SendSysMessage sends a system message to the client
func (c *Client) SendSysMessage(msg string) {
	jsonmessage, _ := json.Marshal(Response{Message: msg, Type: "_SYSMESSAGE"})
	c.send <- jsonmessage
}

// Unregister removes current client from chatAPI
func (c *Client) Unregister(closeChannel bool) {
	var id uint64 = 0
	if c.user != nil {
		id = c.user.ID
	}
	if _, ok := c.cAPI.users[id]; ok {
		if closeChannel {
			close(c.send)
		}
		delete(c.cAPI.users[id], c)
		if len(c.cAPI.users[id]) == 0 {
			delete(c.cAPI.users, id)
		}
	}
}

// Register adds current client connection to chatAPI
func (c *Client) Register() {
	var id uint64 = 0
	if c.user != nil {
		id = c.user.ID
	}
	if _, ok := c.cAPI.users[id]; !ok {
		c.cAPI.users[id] = make(map[*Client]bool)
	}
	c.cAPI.users[id][c] = true
}
