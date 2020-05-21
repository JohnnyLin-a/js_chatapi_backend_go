package chatapi

import (
	"bytes"
	"log"
	"net/http"
	"strconv"
	"time"

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

var (
	userNumber int64 = 1
)

// Client is a bridge between the websocket connection and the ChatAPI.
type Client struct {
	cAPI        *ChatAPI
	conn        *websocket.Conn
	send        chan []byte
	displayName string
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
	// Write to client's html page
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
	log.Println("Origin", r.Header["Origin"][0])
	if r.Header["Origin"][0] != "http://localhost:8080" {
		// TODO: Check for cross-origin resource sharing
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{cAPI: cAPI, conn: conn, send: make(chan []byte, 256), displayName: "User" + strconv.FormatInt(userNumber, 10)}
	log.Println("New user connected: ", client.displayName)
	newUserNumber := &userNumber
	*newUserNumber++
	client.cAPI.register <- client

	go client.startWebsocketReader()
	go client.startWebsocketWriter()
	// send to #general
	go cAPI.handleOnConnect(client)
}
