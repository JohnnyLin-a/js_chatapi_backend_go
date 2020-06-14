package database

import (
	"log"
	"os"

	"github.com/jinzhu/gorm"
	// Blank import for postgres driver
	// _ "github.com/lib/pq"

	// Blank import for postgres ORM dialect
	_ "github.com/jinzhu/gorm/dialects/postgres"

	"github.com/JohnnyLin-a/js_chatapi_backend_go/pkg/chatapi/database/models"
)

// NewDatabase returns a new db connection
func NewDatabase() (*gorm.DB, error) {
	var host = os.Getenv("POSTGRES_HOST")
	var port = os.Getenv("POSTGRES_PORT")
	var user = os.Getenv("POSTGRES_USER")
	var dbname = os.Getenv("POSTGRES_DB")
	var password = os.Getenv("POSTGRES_PASSWORD")
	db, err := gorm.Open("postgres", "sslmode=disable host="+host+" port="+port+" user="+user+" dbname="+dbname+" password="+password)
	if err != nil {
		log.Println("database.NewDatabase: Database connection failed. ", err)
		return nil, err
	}
	return db, nil
}

// Migrate migrates all Models
func Migrate() {
	var db, err = NewDatabase()
	if err != nil {
		return
	}
	defer func() {
		log.Println("Complete.")
		db.Close()
	}()

	// Only migrate Message Modal for now
	log.Println("AutoMigrating...")
	db.AutoMigrate(&models.Message{}, &models.User{}, &models.Chatroom{})
}
