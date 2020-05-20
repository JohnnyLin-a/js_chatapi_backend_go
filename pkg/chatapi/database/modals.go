package database

import (
	"time"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Chatroom struct for db
type Chatroom struct {
	Name            string `json:"name"`
	Mood            string `json:"mood"`
	LastMessage     string `json:"lastMessage"`
	IsDirectMessage bool   `json:"isDirectMessage"`
}

// Message struct for db
type Message struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Chatroom  string    `json:"chatroom"`
	Timestamp time.Time `json:"timestamp"`
	Sender    string    `json:"sender"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
}
