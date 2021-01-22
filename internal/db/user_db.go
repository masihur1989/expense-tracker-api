package db

import (
	"context"
	"log"

	"github.com/masihur1989/expense-tracker-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

// InsertNewUser create a nre record at users collection
func (c *MongoDBClient) InsertNewUser(user *models.User) (interface{}, error) {
	collection := c.Client.Database(c.DBName).Collection("users")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatalf("Error on inserting new User: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// ReadOneUser read a single user
func (c *MongoDBClient) ReadOneUser(filter bson.M) (models.User, error) {
	var user models.User
	collection := c.Client.Database(c.DBName).Collection("users")
	documentReturned := collection.FindOne(context.TODO(), filter)
	documentReturned.Decode(&user)
	return user, nil
}

// ReadAllUsers read all the users
func (c *MongoDBClient) ReadAllUsers(filter interface{}) ([]*models.User, error) {
	var users []*models.User
	collection := c.Client.Database(c.DBName).Collection("users")
	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return users, err
	}
	for cur.Next(context.TODO()) {
		var user models.User
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
func (c *MongoDBClient) RemoveOneUser(filter bson.M) (int64, error) {
	collection := c.Client.Database(c.DBName).Collection("users")
	deleteResult, err := collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal("Error on deleting one Hero", err)
		return 0, err
	}
	return deleteResult.DeletedCount, nil
}

// UpdateOneUser update one user from collections
func (c *MongoDBClient) UpdateOneUser(updatedData bson.M, filter bson.M) (int64, error) {
	collection := c.Client.Database(c.DBName).Collection("users")
	atualizacao := bson.D{{Key: "$set", Value: updatedData}}
	updatedResult, err := collection.UpdateOne(context.TODO(), filter, atualizacao)
	if err != nil {
		log.Fatal("Error on updating one Hero", err)
		return 0, err
	}
	return updatedResult.ModifiedCount, nil
}
