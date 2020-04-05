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
	var host = os.Getenv("PSQL_HOST")
	var port = os.Getenv("PSQL_PORT")
	var user = os.Getenv("PSQL_USER")
	var dbname = os.Getenv("PSQL_DBNAME")
	var password = os.Getenv("PSQL_PASSWORD")
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
	db.Create(m)
	db.Close()
}

// GetLast10Messages gets last 10 messages for a chatroom
func GetLast10Messages(chatroom string) []Message {
	var db, err = newDatabase()
	if err != nil {
		return nil
	}

	var msgs []Message
	db.Order("timestamp desc").Limit(10).Model(&Message{}).Where("chatroom = ?", chatroom).Find(&msgs)

	for _, msg := range msgs {
		log.Println(msg)
	}
	db.Close()
	return msgs
}

// Migrate migrates all Models
func Migrate() {
	var db, err = newDatabase()
	if err != nil {
		return
	}
	log.Println("AutoMigrating...")
	db.AutoMigrate(&Message{})
	db.Close()
}
