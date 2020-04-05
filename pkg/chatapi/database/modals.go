package database

import (
	"time"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Chatroom struct for db
type Chatroom struct {
	Name            string
	Mood            string
	LastMessage     string
	IsDirectMessage bool
}

// Message struct for db
type Message struct {
	ID        string `gorm:"type:uuid;primary_key;default:uuid_generate_v4()"`
	Chatroom  string
	Timestamp time.Time
	Sender    string
	Type      string
	Message   string
}
