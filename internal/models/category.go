package models

import (
	"context"
	"log"
	"time"

	"github.com/masihur1989/expense-tracker-api/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category model for category collection
type Category struct {
	ID        primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
	Name      string             `json:"name" bson:"name" validate:"required,alpha"`
}

// CategoryUpdateInput model for update endpoint
type CategoryUpdateInput struct {
	Name string `json:"name" bson:"name" validate:"required,alpha"`
}

// CategoryModeler godoc
type CategoryModeler interface {
	Insert(catergory *Category) (interface{}, error)
	ReadAll(filter interface{}) ([]Category, error)
	ReadOne(filter interface{}) (Category, error)
	UpdateOne(updatedData interface{}, filter interface{}) (int64, error)
	RemoveOne(filter interface{}) (int64, error)
}

// CategoryModel godoc
type CategoryModel struct {
	db db.MongoDBClient
}

// NewCategoryModel godoc
func NewCategoryModel(db db.MongoDBClient) *CategoryModel {
	return &CategoryModel{db}
}

// Insert insert a record at categories collection
func (c *CategoryModel) Insert(catergory *Category) (interface{}, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("categories")
	insertResult, err := collection.InsertOne(context.TODO(), catergory)
	if err != nil {
		log.Fatalf("Error on inserting new category: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// ReadAll read all the categories
func (c *CategoryModel) ReadAll(filter interface{}) ([]Category, error) {
	var categories []Category
	collection := c.db.Client.Database(c.db.DBName).Collection("categories")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return categories, err
	}
	for cur.Next(context.TODO()) {
		var category Category
		err = cur.Decode(&category)
		if err != nil {
			log.Printf("Error on Decoding the document: %v\n", err)
		}
		categories = append(categories, category)
	}
	log.Printf("documentReturned: %v\n", categories)
	return categories, nil
}

// ReadOne read a single category
func (c *CategoryModel) ReadOne(filter interface{}) (Category, error) {
	var category Category
	collection := c.db.Client.Database(c.db.DBName).Collection("categories")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&category)
	return category, nil
}

// UpdateOne update one category from collections
func (c *CategoryModel) UpdateOne(updatedData interface{}, filter interface{}) (int64, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("categories")
	atualizacao := bson.D{{Key: "$set", Value: updatedData}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
		log.Fatal("Error on updating one Hero", err)
		return 0, err
	}
	return updatedResult.ModifiedCount, nil
}

// RemoveOne remove one category from collections
func (c *CategoryModel) RemoveOne(filter interface{}) (int64, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("categories")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleting one Category", err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}
