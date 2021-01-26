package models

import (
	"context"
	"log"
	"time"

	"github.com/masihur1989/expense-tracker-api/internal/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// ProjectDetails collection structure
type ProjectDetails struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
	Expenses    []Expense          `json:"expenses" bson:"expenses"`
	Users       []User             `json:"users" bson:"users"`
}

// ProjectExpenseQS Query String parser for project-details query
type ProjectExpenseQS struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
}

// Project collection structure
type Project struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	Title       string             `json:"title" bson:"title"`
	Description string             `json:"description" bson:"description"`
}

// ProjectUser collection structure for projectUser
type ProjectUser struct {
	ID          primitive.ObjectID `json:"id" bson:"_id"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	ProjectID   primitive.ObjectID `json:"project_id" bson:"project_id"`
	Email       string             `json:"email" bson:"email" validate:"required,email"`                 // TODO needs to decide of if we want to make it unique
	PhoneNumber string             `json:"phone_number" bson:"phone_number" validate:"required,numeric"` // TODO needs to decide of if we want to make it unique
	Name        string             `json:"name" bson:"name" validate:"required,alpha"`
	Role        Role               `json:"role" bson:"role" validate:"required,oneof=ADMIN SUPERVISOR STAFF USER"`
	IsActive    bool               `json:"is_active" bson:"is_active" validate:"required"`
}

// ProjectModeler godoc
type ProjectModeler interface {
	Insert(project *Project) (interface{}, error)
	ReadAll(filter interface{}) ([]Project, error)
	ReadOne(filter interface{}) (Project, error)
	LookupProjectDetails(filter bson.D, expfilter ProjectExpenseQS) (ProjectDetails, error)
	InsertProjectUser(projectUser *ProjectUser) (interface{}, error)
}

// ProjectModel godoc
type ProjectModel struct {
	db db.MongoDBClient
}

// NewProjectModel godoc
func NewProjectModel(db db.MongoDBClient) *ProjectModel {
	return &ProjectModel{db}
}

// Insert insert a record at projects collection
func (c *ProjectModel) Insert(project *Project) (interface{}, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("projects")
	insertResult, err := collection.InsertOne(context.TODO(), project)
	if err != nil {
		log.Fatalf("Error on inserting new project: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}

// ReadAll read all the projects
func (c *ProjectModel) ReadAll(filter interface{}) ([]Project, error) {
	var projects []Project
	collection := c.db.Client.Database(c.db.DBName).Collection("projects")

	cur, err := collection.Find(context.TODO(), filter)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return projects, err
	}
	for cur.Next(context.TODO()) {
		var project Project

		err = cur.Decode(&project)
		if err != nil {
			log.Printf("Error on Decoding the document: %v\n", err)
		}

		projects = append(projects, project)
	}
	log.Printf("documentReturned: %v\n", projects)
	return projects, nil
}

// ReadOne read a single project
func (c *ProjectModel) ReadOne(filter interface{}) (Project, error) {
	var project Project
	collection := c.db.Client.Database(c.db.DBName).Collection("projects")
	projectReturned := collection.FindOne(context.TODO(), filter)
	projectReturned.Decode(&project)
	return project, nil
}

// LookupProjectDetails parse all the project details with the project_id
func (c *ProjectModel) LookupProjectDetails(projectFilter bson.D, expfilter ProjectExpenseQS) (ProjectDetails, error) {
	var project ProjectDetails
	collection := c.db.Client.Database(c.db.DBName).Collection("projects")
	pipeline := mongo.Pipeline{
		{{"$match", projectFilter}},
		{{"$lookup", bson.D{
			{"from", "expenses"},
			{"localField", "_id"},
			{"foreignField", "project_id"},
			{"as", "expenses"},
		}}},
		{{"$sort", bson.D{
			{"date", -1},
		}}},
		{{"$lookup", bson.D{
			{"from", "projectUsers"},
			{"localField", "_id"},
			{"foreignField", "project_id"},
			{"as", "users"},
		}}},
		{{"$project", bson.D{
			{"_id", 1},
			{"title", 1},
			{"description", 1},
			{"created_at", 1},
			{"updated_at", 1},
			{"expenses", bson.D{
				{"$filter", bson.D{
					{"input", "$expenses"},
					{"as", "expense"},
					{"cond", bson.D{
						{"$and", bson.A{
							bson.D{{"$gte", bson.A{"$$expense.date", expfilter.Start}}},
							bson.D{{"$lt", bson.A{"$$expense.date", expfilter.End}}},
						}},
					}},
				}},
			}},
			{"users", 1},
		}}},
	}

	cur, err := collection.Aggregate(context.TODO(), pipeline)
	if err != nil {
		log.Printf("ERROR FINDING DATA: %v\n", err)
		return project, err
	}
	for cur.Next(context.TODO()) {
		err = cur.Decode(&project)
		if err != nil {
			log.Printf("Error on Decoding the document: %v\n", err)
			return project, err
		}
	}

	return project, nil
}

// InsertProjectUser insert a record at projectUsers collection
func (c *ProjectModel) InsertProjectUser(user *ProjectUser) (interface{}, error) {
	collection := c.db.Client.Database(c.db.DBName).Collection("projectUsers")
	insertResult, err := collection.InsertOne(context.TODO(), user)
	if err != nil {
		log.Fatalf("Error on inserting new project user: %v\n", err)
		return nil, err
	}
	return insertResult.InsertedID, nil
}
