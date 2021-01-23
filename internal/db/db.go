package db

import (
	"context"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDBClient mongo db client needed for the project
type MongoDBClient struct {
	Client *mongo.Client
	DBName string
}

// Used to create a singleton object of MongoDB client.
// Initialized and exposed through  GetMongoClient().
var clientInstance *mongo.Client

//Used during creation of singleton client object in GetMongoClient().
var clientInstanceError error

//Used to execute client creation procedure only once.
var mongoOnce sync.Once

var mongoDbInstance, dbInstance string

// GetClient get the db client
func GetClient() (MongoDBClient, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
		mongoDbInstance = utils.MustGet("MONGO_DB_INSTANCE")
		dbInstance = utils.MustGet("DB_INSTANCE")
		// Set client options
		clientOptions := options.Client().ApplyURI(mongoDbInstance)
		// Connect to MongoDB
		client, err := mongo.Connect(context.TODO(), clientOptions)
		if err != nil {
			clientInstanceError = err
		}
		// Check the connection
		err = client.Ping(context.TODO(), nil)
		if err != nil {
			clientInstanceError = err
		}
		clientInstance = client
	})
	return MongoDBClient{
		Client: clientInstance,
		DBName: dbInstance,
	}, clientInstanceError
}
