package models

import (
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for user to map mongodb document
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Email       string             `json:"email" bson:"email" validate:"required"`               // TODO needs to decide of if we want to make it unique
	PhoneNumber string             `json:"phone_number" bson:"phone_number" validate:"required"` // TODO needs to decide of if we want to make it unique
	Name        string             `json:"name" bson:"name" validate:"required"`
	Role        Role               `json:"role" bson:"role" validate:"required"`
	IsActive    bool               `json:"is_active" bson:"is_active" validate:"required"`
}

// UserUpdateInput godoc
type UserUpdateInput struct {
	Name     string `json:"name" bson:"name" validate:"required,max=20" `
	IsActive bool   `json:"is_active" bson:"is_active" validate:"required"`
}

// Role user role
type Role string

// all the acceptable roles
const (
	RoleAdmin      Role = "ADMIN"
	RoleSupervisor Role = "SUPERVISOR"
	RoleStaff      Role = "STAFF"
	RoleUser       Role = "USER"
)

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
