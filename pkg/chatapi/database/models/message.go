package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Message struct for db
type Message struct {
	ID         string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	ChatroomID string    `gorm:"foreignkey:ChatroomRefer" json:"chatroomID"`
	UserID     uint64    `gorm:"foreignkey:UserRefer" json:"userID"`
	Timestamp  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"timestamp"`
	Type       string    `gorm:"not null" json:"type"`
	Message    string    `gorm:"not null" json:"message"`
}

// Save saves a message Model to the database
func (m *Message) Save(db *gorm.DB) error {
	if db == nil {
		return errors.New("No DB connection")
	}
	return db.Create(m).Error
}

// GetLast100Messages gets last 10 messages for a chatroom
func GetLast100Messages(db *gorm.DB, chatroom *string) (*[]Message, error) {
	var msgs []Message
	if db == nil {
		return nil, errors.New("No DB connection")
	}

	return &msgs, db.Raw("SELECT * FROM (SELECT * FROM messages WHERE chatroom_id = ? ORDER BY timestamp DESC LIMIT 100) AS sq ORDER BY sq.timestamp ASC;", *chatroom).Scan(&msgs).Error
}
