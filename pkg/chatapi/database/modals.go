package database

import (
	"time"

	"github.com/jinzhu/gorm"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Chatroom struct for db
type Chatroom struct {
	gorm.Model
	Name            string
	Mood            string
	LastMessage     string
	IsDirectMessage bool
}

// Message struct for db
type Message struct {
	Chatroom  string
	Timestamp time.Time
	Sender    string
	Type      string
	Message   string
}
