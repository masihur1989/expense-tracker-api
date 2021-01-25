package handler

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// objectIDFromStringID convert the string id to primitive.ObjectID
func objectIDFromStringID(param string) (primitive.ObjectID, error) {
	// need to convert it to ObjectID.
	id, err := primitive.ObjectIDFromHex(param)
	if err != nil {
		log.Printf("INVALID ID PASSED: %v\n", err)
		return primitive.NilObjectID, err
	}
	return id, nil
}

// parseDateToFormat parse date string to desired fromat
func parseDateToFormat(layout, date string) (time.Time, error) {
	return time.Parse(layout, date)
}
