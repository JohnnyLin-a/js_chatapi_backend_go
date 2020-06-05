package models

import (
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