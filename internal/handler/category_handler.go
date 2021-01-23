package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryHandler godoc
type CategoryHandler struct {
	catModel models.CategoryModeler
}

// NewCategoryHandler godoc
func NewCategoryHandler(cm models.CategoryModeler) CategoryHandler {
	return CategoryHandler{cm}
}

// CreateCategory godoc
// @Summary Create category.
// @Description create category.
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Create Category"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/categories [post]
func (c CategoryHandler) CreateCategory(e echo.Context) error {
	cat := new(models.Category)
	if err := e.Bind(cat); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := e.Validate(cat); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	// fill the nil values
	cat.ID = primitive.NewObjectID()
	cat.CreatedAt = time.Now()
	cat.UpdatedAt = time.Now()

	id, err := c.catModel.Insert(cat)

	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}

	return utils.Data(http.StatusCreated, id, "category created", e)
}

// GetCategories godoc
// get all the users
// QueryParams accepted are name
// and only single value is used. TODO needs to find a way to accept array of query pamra values
// @Summary Get Categories.
// @Description get category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param name query string false "name search by name"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/categories [get]
func (c CategoryHandler) GetCategories(e echo.Context) error {
	qs := e.QueryParams()
	var filter interface{}
	if len(e.QueryParams()) == 0 {
		filter = bson.D{}
	} else {
		f := bson.M{}
		if x, ok := qs["name"]; ok {
			f["name"] = x[0]
		}
		filter = f
	}

	cats, err := c.catModel.ReadAll(filter)
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}
	return utils.Data(http.StatusOK, cats, "category details", e)
}

// DeleteCategory godoc
// @Summary Delete a Category.
// @Description get category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/categories/{id} [delete]
func (c CategoryHandler) DeleteCategory(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	count, err := c.catModel.RemoveOne(bson.M{"_id": ID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusAccepted, count, "category removed", e)
}

// UpdateCategory godoc
// @Summary Update a Category.
// @Description update category by ID
// @Tags categories
// @Accept json
// @Produce json
// @Param id path string true "Category ID"
// @Param category body models.CategoryUpdateInput true "Update Category"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/categories/{id} [put]
func (c CategoryHandler) UpdateCategory(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	catInput := new(models.CategoryUpdateInput)

	if err := e.Bind(catInput); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := e.Validate(catInput); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	// update fields - name
	update := bson.M{
		"name": catInput.Name,
	}

	count, err := c.catModel.UpdateOne(update, bson.M{"_id": ID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusOK, count, "category updated", e)
}
