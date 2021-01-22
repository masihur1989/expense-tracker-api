package db

import (
	"context"
	"sync"

	"github.com/masihur1989/expense-tracker-api/internal/models"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Client Interface to hold all the client related func
type Client interface {
	InsertNewUser(user *models.User) (interface{}, error)
	ReadOneUser(filter bson.M) (models.User, error)
	ReadAllUsers(filter interface{}) ([]*models.User, error)
	RemoveOneUser(filter bson.M) (int64, error)
}

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

func init() {
	mongoDbInstance = utils.MustGet("MONGO_DB_INSTANCE")
	dbInstance = utils.MustGet("DB_INSTANCE")
}

// GetClient get the db client
func GetClient() (MongoDBClient, error) {
	//Perform connection creation operation only once.
	mongoOnce.Do(func() {
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
