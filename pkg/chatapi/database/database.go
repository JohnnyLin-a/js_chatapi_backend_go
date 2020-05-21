package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// Blank import for postgres driver
	// _ "github.com/lib/pq"

	// Blank import for postgres ORM dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func newDatabase() (*gorm.DB, error) {
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var dbname = os.Getenv("POSTGRES_DB")
	var password = os.Getenv("POSTGRES_PASSWORD")
	db, err := gorm.Open("postgres", "sslmode=disable host="+host+" port="+port+" user="+user+" dbname="+dbname+" password="+password)
	if err != nil {
		log.Println("database.newDatabase: Database connection failed. ", err)
		return nil, err
	}
	return db, nil
}

// SaveMessage saves a message Model to the database
func SaveMessage(m *Message) {
	var db, err = newDatabase()
	if err != nil {
		return
	}
	defer func() {
		db.Close()
	}()
	db.Create(m)
}

// GetLast100Messages gets last 10 messages for a chatroom
func GetLast100Messages(chatroom *string) []Message {
	var db, err = newDatabase()
	if err != nil {
		return nil
	}
	defer func() {
		db.Close()
	}()

	var msgs []Message
	db.Raw("SELECT * FROM (SELECT * FROM messages WHERE chatroom = ? ORDER BY timestamp DESC LIMIT 100) AS sq ORDER BY sq.timestamp ASC;", *chatroom).Scan(&msgs)
	return msgs
}

// Migrate migrates all Models
func Migrate() {
	var db, err = newDatabase()
	if err != nil {
		return
	}
	defer func() {
		log.Println("Complete.")
		db.Close()
	}()

	// Only migrate Message Modal for now
	log.Println("AutoMigrating...")
	db.AutoMigrate(&Message{})
}
