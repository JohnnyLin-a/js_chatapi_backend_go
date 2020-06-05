package models

import (
	"time"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/jinzhu/gorm"
)

// Message struct for db
type Message struct {
	ID        string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"-"`
	Chatroom  string    `gorm:"default:'#general'" json:"chatroom"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Sender    string    `json:"sender"`
	Type      string    `json:"type"`
	Message   string    `json:"message"`
}


// SaveMessage saves a message Model to the database
func (m *Message) SaveMessage(db *gorm.DB) {
	if db == nil {
		return
	}
	defer func() {
		db.Close()
	}()
	m.Timestamp = time.Now()
	db.Create(m)
}

// GetLast100Messages gets last 10 messages for a chatroom
func GetLast100Messages(db *gorm.DB, chatroom *string) []Message {
	var msgs []Message
	if db == nil {
		return msgs
	}
	defer func() {
		db.Close()
	}()
	
	db.Raw("SELECT * FROM (SELECT * FROM messages WHERE chatroom = ? ORDER BY timestamp DESC LIMIT 100) AS sq ORDER BY sq.timestamp ASC;", *chatroom).Scan(&msgs)
	return msgs
}