package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User is the user model that is used for the graphql queries and database tables
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email" gorm:"unique"`
	Username  string    `json:"username" gorm:"unique"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New().String()

	return nil
}
