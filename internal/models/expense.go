package models

import (
	"context"
	"log"
	"time"

	"github.com/masihur1989/expense-tracker-api/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Expense expesne model
type Expense struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Date        time.Time          `json:"date" bson:"date"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Location    string             `json:"location" bson:"location"`
	Total       float64            `json:"total" bson:"total"`
	Status      string             `json:"status" bson:"status"`
	ProjectID   primitive.ObjectID `json:"project_id" bson:"project_id"`
	Category    Category           `json:"category" bson:"category"`
	InsertedBy  User               `json:"user" bson:"user"`
}

// ExpenseInput expense create input model
type ExpenseInput struct {
	Date        string  `json:"date" bson:"date" validate:"required"` // string date give more controll to parse it in any form for storage
	Title       string  `json:"title" bson:"title" validate:"required"`
	Description string  `json:"description" bson:"description" validate:"required"`
	Location    string  `json:"location" bson:"location"`
	Total       float64 `json:"total" bson:"total" validate:"required"`
	Status      string  `json:"status" bson:"status" validate:"required,oneof=pending confirmed"`
	CategoryID  string  `json:"category_id" bson:"category_id" validate:"required"`
	InsertedBy  string  `json:"inserted_by" bson:"inserted_by" validate:"required"`
}

// ExpenseModeler godoc
type ExpenseModeler interface {
	Insert(expense Expense) (interface{}, error)
	ReadAll(filter interface{}) ([]Expense, error)
	ReadOne(filter interface{}) (Expense, error)
	Remove(filter interface{}) (int64, error)
	UpdateOne(updatedData interface{}, filter interface{}) (int64, error)
}

// ExpenseModel godoc
type ExpenseModel struct {
	db db.MongoDBClient
}

// NewExpenseModel godoc
func NewExpenseModel(db db.MongoDBClient) *ExpenseModel {
	return &ExpenseModel{db}
}

// Insert insert a record at expenses collection
func (e *ExpenseModel) Insert(expense Expense) (interface{}, error) {
	collection := e.db.Client.Database(e.db.DBName).Collection("expenses")
	insertResult, err := collection.InsertOne(context.TODO(), expense)
	if err != nil {
		log.Fatalf("Error on inserting new expense: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// ReadAll read all the expenses
func (e *ExpenseModel) ReadAll(filter interface{}) ([]Expense, error) {
	var expenses []Expense
	collection := e.db.Client.Database(e.db.DBName).Collection("expenses")
	log.Printf("filter: %v\n", filter)
	// sort the entries based on the `date` field
	opts := options.FindOptions{}
	opts.SetSort(bson.D{{"date", -1}})
	cur, err := collection.Find(context.TODO(), filter, &opts)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return expenses, err
	}
	for cur.Next(context.TODO()) {
		var expense Expense
		err = cur.Decode(&expense)
		if err != nil {
			log.Printf("Error on Decoding the document: %v\n", err)
		}
		expenses = append(expenses, expense)
	}
	log.Printf("documentReturned: %v\n", expenses)
	return expenses, nil
}

// ReadOne read a single expense
func (e *ExpenseModel) ReadOne(filter interface{}) (Expense, error) {
	var expense Expense
	collection := e.db.Client.Database(e.db.DBName).Collection("expenses")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&expense)
	return expense, nil
}

// Remove remove one expense from collctions
func (e *ExpenseModel) Remove(filter interface{}) (int64, error) {
	collection := e.db.Client.Database(e.db.DBName).Collection("expenses")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleting one expense", err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

// UpdateOne update one expense from collections
func (e *ExpenseModel) UpdateOne(updatedData interface{}, filter interface{}) (int64, error) {
	collection := e.db.Client.Database(e.db.DBName).Collection("expenses")
	atualizacao := bson.D{{Key: "$set", Value: updatedData}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
		log.Fatal("Error on updating one expense", err)
		return 0, err
	}
	return updatedResult.ModifiedCount, nil
}
