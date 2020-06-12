package models

import (
	"time"

	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"

	// blank import for postgres
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

// User struct for database storage
type User struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"-"`
	CreatedAt   time.Time `json:"-"`
	DeletedAt   time.Time `json:"-"`
	UpdatedAt   time.Time `json:"-"`
	DisplayName string    `gorm:"size:255;not null" json:"displayName"`
	Email       string    `gorm:"size:255;unique;not null" json:"email"`
	Username    string    `gorm:"size:255;unique;not null" json:"username"`
	Password    string    `gorm:"not null" json:"-"`
}

// Hash hashes a password string
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword verifies a password string against a hashed password
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// DeleteUser soft deletes a user.
func (u *User) DeleteUser(db *gorm.DB) {
	if db == nil {
		return
	}
	defer func() {
		db.Close()
	}()
	db.Delete(u)
}

// SaveUser saves a new user
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	if err := db.Create(&u).Error; err != nil {
		return &User{}, err
	}
	return u, nil
}

// FindUserByID tries to find a user by ID
func FindUserByID(db *gorm.DB, id uint64) (*User, error) {
	var u User
	if err := db.First(&u, id).Error; err != nil {
		return &User{}, err
	}
	return &u, nil
}
