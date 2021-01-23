package handler

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/labstack/echo/v4"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserModelStub struct{}

var obzID primitive.ObjectID
var err error

func init() {
	godotenv.Load("../../.env")
	obzID, err = primitive.ObjectIDFromHex("6009be17d6a899ab8340eb79")
	if err != nil {
		log.Printf("ERROR %v", err)
	}
}

func (u UserModelStub) ReadOneUser(filter interface{}) (models.User, error) {
	return models.User{
		ID:          obzID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Email:       "marufrahman1349@gmail.com",
		PhoneNumber: "01234567890",
		Name:        "test",
		Role:        "USER",
		IsActive:    true,
	}, nil
}

func (u UserModelStub) InsertNewUser(user *models.User) (interface{}, error) {
	return 0, nil
}

func (u UserModelStub) ReadAllUsers(filter interface{}) ([]*models.User, error) {
	var users []*models.User
	users = append(users, &models.User{
		ID:          obzID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Email:       "test.user@gmail.com",
		PhoneNumber: "01234567890",
		Name:        "test",
		Role:        "USER",
		IsActive:    true,
	})
	return users, nil
}

func (u UserModelStub) RemoveOneUser(filter interface{}) (int64, error) {
	return 0, nil
}

func (u UserModelStub) UpdateOneUser(updatedData interface{}, filter interface{}) (int64, error) {
	return 0, nil
}

func TestReadOneUser(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")
	c.SetParamNames("id")
	c.SetParamValues("6009be17d6a899ab8340eb79")

	u := UserModelStub{}
	h := NewUserHandler(u)

	if assert.NoError(t, h.GetUser(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestReadAllUsers(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/:id")

	u := UserModelStub{}
	h := NewUserHandler(u)

	if assert.NoError(t, h.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func TestInsertNewUser(t *testing.T) {
	e := echo.New()
	reqByte, err := json.Marshal(models.User{
		ID:          obzID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Email:       "test.user@gmail.com",
		PhoneNumber: "01234567890",
		Name:        "test",
		Role:        "USER",
		IsActive:    true,
	})

	if err != nil {
		log.Printf("ERROR MARSHALING STRUCT: %v", err)
	}
	requestReader := bytes.NewReader(reqByte)
	req := httptest.NewRequest(echo.POST, "/", requestReader)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/api/v1/users/")

	u := UserModelStub{}
	h := NewUserHandler(u)

	if assert.NoError(t, h.GetUsers(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}
