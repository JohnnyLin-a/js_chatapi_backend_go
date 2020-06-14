package models

import (
	"errors"

	"github.com/jinzhu/gorm"
	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// Chatroom struct for db model
type Chatroom struct {
	ID              string    `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name            *string   `json:"name"`
	Mood            *string   `json:"mood"`
	LastMessage     *string   `json:"lastMessage"`
	IsDirectMessage bool      `json:"isDirectMessage"`
	Users           []User    `json:"-"`
	Messages        []Message `json:"-"`
}

// Save saves a Chatroom
func (cr *Chatroom) Save(db *gorm.DB) (err error) {
	if db == nil {
		return errors.New("DB Connection error")
	}
	return db.Create(cr).Error
}
