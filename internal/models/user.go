package models

import (
	"net/url"
	"regexp"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User model for user to map mongodb document
type User struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Email       string             `json:"email" bson:"email"`               // TODO needs to decide of if we want to make it unique
	PhoneNumber string             `json:"phone_number" bson:"phone_number"` // TODO needs to decide of if we want to make it unique
	Name        string             `json:"name" bson:"name"`
	Role        Role               `json:"role" bson:"role"`
	IsActive    bool               `json:"is_active" bson:"is_active"`
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

// UserPostValidator validated the user creation input
func (user *User) UserPostValidator() url.Values {
	errs := url.Values{}

	if user.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	if user.PhoneNumber == "" {
		errs.Add("phone number", "The Phone Number field is required!")
	}

	if user.Email == "" {
		errs.Add("email", "The Email field is required!")
	}

	if !isEmailValid(user.Email) {
		errs.Add("email", "The Email field is invalid!")
	}

	if user.Role == "" {
		errs.Add("role", "The Role field is required!")
	}

	return errs
}

// UserPutValidator validated the user creation input
func (user *User) UserPutValidator() url.Values {
	errs := url.Values{}
	if user.Name == "" {
		errs.Add("name", "The name field is required!")
	}

	return errs
}

// isEmailValid checks if the email provided passes the required structure and length.
func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
