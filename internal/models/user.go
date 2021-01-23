package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for user to map mongodb document
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`                 // TODO needs to decide of if we want to make it unique
	PhoneNumber string             `json:"phone_number" bson:"phone_number" validate:"required,numeric"` // TODO needs to decide of if we want to make it unique
	Name        string             `json:"name" bson:"name" validate:"required,alpha"`
	Role        Role               `json:"role" bson:"role" validate:"required,oneof=ADMIN SUPERVISOR STAFF USER"`
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
