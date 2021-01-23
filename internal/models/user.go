package models

import (
	"context"
	"log"
	"time"

	"github.com/masihur1989/expense-tracker-api/internal/db"
	"go.mongodb.org/mongo-driver/bson"
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

// UserModel godoc
type UserModel interface {
	InsertNewUser(user *User) (interface{}, error)
	ReadOneUser(filter interface{}) (User, error)
	ReadAllUsers(filter interface{}) ([]*User, error)
	RemoveOneUser(filter interface{}) (int64, error)
	UpdateOneUser(updatedData interface{}, filter interface{}) (int64, error)
}

// UserModelImpl godoc
type UserModelImpl struct {
	db db.MongoDBClient
}

// NewUserModelImpl godoc
func NewUserModelImpl(db db.MongoDBClient) *UserModelImpl {
	return &UserModelImpl{
		db: db,
	}
}

// InsertNewUser create a nre record at users collection
func (c *UserModelImpl) InsertNewUser(user *User) (interface{}, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatalf("Error on inserting new User: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// ReadOneUser read a single user
func (c *UserModelImpl) ReadOneUser(filter interface{}) (User, error) {
	var user User
	collection := c.db.Client.Database(c.db.DBName).Collection("users")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&user)
	return user, nil
}

// ReadAllUsers read all the users
func (c *UserModelImpl) ReadAllUsers(filter interface{}) ([]*User, error) {
	var users []*User
	collection := c.db.Client.Database(c.db.DBName).Collection("users")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return users, err
	}
	for cur.Next(context.TODO()) {
		var user User
		err = cur.Decode(&user)
		if err != nil {
			log.Printf("Error on Decoding the document: %v\n", err)
		}
		users = append(users, &user)
	}
	log.Printf("documentReturned: %v\n", users)
	return users, nil
}

// RemoveOneUser remove one user from collctions
func (c *UserModelImpl) RemoveOneUser(filter interface{}) (int64, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("users")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleting one Hero", err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

// UpdateOneUser update one user from collections
func (c *UserModelImpl) UpdateOneUser(updatedData interface{}, filter interface{}) (int64, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("users")
	atualizacao := bson.D{{Key: "$set", Value: updatedData}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
		log.Fatal("Error on updating one Hero", err)
		return 0, err
	}
	return updatedResult.ModifiedCount, nil
}
