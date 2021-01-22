package handler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/masihur1989/expense-tracker-api/internal/db"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserHandler controller for users
type UserHandler struct {
	DBClient db.MongoDBClient
}

// NewUserHandler scho.Echo handler function
func NewUserHandler(e *echo.Group, dbClient db.MongoDBClient) {
	userController := UserHandler{
		DBClient: dbClient,
	}
	// users routes
	e.GET("/users/", userController.GetUsers)
	e.GET("/users/:id", userController.GetUser)
	e.POST("/users", userController.CreateUser)
	e.PUT("/users/:id", userController.UpdateUser)
	e.DELETE("/users/:id", userController.DeleteUser)
}

// CreateUser godoc
// @Summary Create User.
// @Description create user.
// @Tags users
// @Accept json
// @Produce json
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/users [post]
func (u *UserHandler) CreateUser(c echo.Context) error {
	user := new(models.User)
	// fill the nill values
	user.ID = primitive.NewObjectID()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := c.Bind(user); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if validErrs := user.UserPostValidator(); len(validErrs) > 0 {
		err := map[string]interface{}{"validationError": validErrs}
		return c.JSON(http.StatusBadRequest, err)
	}

	id, err := u.DBClient.InsertNewUser(user)

	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), c)
	}

	return utils.Data(http.StatusCreated, id, "user created", c)
}

// GetUsers godoc
// get all the users
// QueryParams accepted are is_active, role
// and only single value is used. TODO needs to find a way to accept array of query pamra values
// @Summary Get Users.
// @Description get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param name query string false "name search by name"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/users [get]
func (u *UserHandler) GetUsers(c echo.Context) error {
	qs := c.QueryParams()
	var filter interface{}
	if len(c.QueryParams()) == 0 {
		filter = bson.D{}
	} else {
		f := bson.M{}
		if x, ok := qs["is_active"]; ok {
			b, err := strconv.ParseBool(x[0])
			if err != nil {
				log.Printf("INVALID QUERY PARAM PASSED: %v\n", err)
				return utils.Error(http.StatusBadRequest, err.Error(), c)
			}
			f["is_active"] = b
		}
		if x, ok := qs["role"]; ok {
			f["role"] = x[0]
		}

		filter = f
	}

	users, err := u.DBClient.ReadAllUsers(filter)
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusInternalServerError, err.Error(), c)
	}
	return utils.Data(http.StatusOK, users, "user details", c)
}

// GetUser godoc
// @Summary Get an User.
// @Description get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/users/{id} [get]
func (u *UserHandler) GetUser(c echo.Context) error {
	userID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}
	user, err := u.DBClient.ReadOneUser(bson.M{"_id": userID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusOK, user, "user detail", c)
}

// DeleteUser godoc
// @Summary Delete an User.
// @Description get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/{id} [delete]
func (u *UserHandler) DeleteUser(c echo.Context) error {
	userID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	count, err := u.DBClient.RemoveOneUser(bson.M{"_id": userID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusAccepted, count, "user removed", c)
}

// UpdateUser godoc
// @Summary Update an User.
// @Description update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/users/{id} [put]
func (u *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	user := new(models.User)

	if err := c.Bind(user); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	// update fields - name, is_active, updated_at
	update := bson.M{
		"name":       user.Name,
		"is_active":  user.IsActive,
		"updated_at": time.Now(),
	}

	count, err := u.DBClient.UpdateOneUser(update, bson.M{"_id": userID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusOK, count, "user updated", c)
}
