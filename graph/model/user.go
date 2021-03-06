package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/sethigeet/gql-go-auth-backend/util"
)

// User is the user model that is used for the graphql queries and database tables
type User struct {
	ID                   string    `json:"id" gorm:"type:uuid;primaryKey"`
	Email                string    `json:"email" gorm:"not null;unique;size:256"`
	Username             string    `json:"username" gorm:"not null;unique;size:256"`
	Password             string    `json:"password" gorm:"not null;size:256"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
	Confirmed            bool      `json:"confirmed" gorm:"type:bool;not null;default:false"`
	ForgotPasswordLocked bool      `json:"forgotPasswordLocked" gorm:"type:bool;not null;default:false"`
}

// BeforeCreate is run just before inserting a new user in the users table
// It create a unique id for the user and hashes the password of the user
func (user *User) BeforeCreate(tx *gorm.DB) error {
	user.ID = uuid.New().String()

	hashedPasswd, err := util.HashPassword(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPasswd

	return nil
}
