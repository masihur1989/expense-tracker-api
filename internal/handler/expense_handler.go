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

// ExpenseHandler godoc
type ExpenseHandler struct {
	expenseModel  models.ExpenseModeler
	userModel     models.UserModel
	categoryModel models.CategoryModeler
}

// NewExpenseHandler godoc
func NewExpenseHandler(em models.ExpenseModeler, um models.UserModel, cm models.CategoryModeler) ExpenseHandler {
	return ExpenseHandler{em, um, cm}
}

// CreateExpense godoc
// @Summary Create expense.
// @Description create expense.
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body models.ExpenseInput true "Create Expense"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/expenses [post]
func (e ExpenseHandler) CreateExpense(c echo.Context) error {
	expInput := new(models.ExpenseInput)
	if err := c.Bind(expInput); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := c.Validate(expInput); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	categoryID, err := objectIDFromStringID(expInput.CategoryID)
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	category, err := e.categoryModel.ReadOne(bson.M{"_id": categoryID})
	if err != nil {
		log.Printf("CATEGORY NOT FOUND ERROR: %v\n", err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}

	userID, err := objectIDFromStringID(expInput.InsertedBy)
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	user, err := e.userModel.ReadOneUser(bson.M{"_id": userID})
	if err != nil {
		log.Printf("USER NOT FOUND ERROR: %v\n", err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}

	d, err := parseDateToFormat("2006-01-02", expInput.Date)
	if err != nil {
		log.Printf("Time Parsing Error %v", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), c)
	}

	exp := models.Expense{
		ID:          primitive.NewObjectID(),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Title:       expInput.Title,
		Description: expInput.Description,
		Date:        d,
		Category:    category,
		Location:    expInput.Location,
		Total:       expInput.Total,
		Status:      expInput.Status,
		InsertedBy:  user,
	}

	id, err := e.expenseModel.Insert(exp)

	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), c)
	}

	return utils.Data(http.StatusCreated, id, "expense created", c)
}

// GetExpenses godoc
// get all the expenses
// QueryParams accepted are name
// and only single value is used. TODO needs to find a way to accept array of query pamra values
// @Summary Get Expenses.
// @Description get expenses
// @Tags expenses
// @Accept json
// @Produce json
// @Param name query string false "date search by date"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/expenses [get]
func (e ExpenseHandler) GetExpenses(c echo.Context) error {
	qs := c.QueryParams()
	log.Printf("QS: %s\n", qs)
	var filter interface{}
	if len(c.QueryParams()) == 0 {
		filter = bson.D{}
	} else {
		f := make(map[string]string)
		if x, ok := qs["start"]; ok {
			f["start"] = x[0]
		}
		if x, ok := qs["end"]; ok {
			f["end"] = x[0]
		}
		log.Printf("f %v\n", f)
		var startDate, endDate time.Time
		var err error

		startDate, err = parseDateToFormat("2006-01-02", f["start"])
		if err != nil {
			log.Printf("ERROR PARSING STARTDATE: %v\n", err)
			return utils.Error(http.StatusBadRequest, err.Error(), c)
		}

		endDate, err = parseDateToFormat("2006-01-02", f["end"])
		if err != nil {
			log.Printf("ERROR PARSING ENDDATE: %v\n", err)
			return utils.Error(http.StatusBadRequest, err.Error(), c)
		}

		filter = bson.D{
			{"date", bson.D{
				{"$gte", startDate},
				{"$lt", endDate}, // time.Now().Add(-1 * 24 * time.Hour)
			}},
		}
	}

	cats, err := e.expenseModel.ReadAll(filter)
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusInternalServerError, err.Error(), c)
	}
	return utils.Data(http.StatusOK, cats, "expense details", c)
}

// GetExpense godoc
// @Summary Get an Expense.
// @Description get expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/expenses/{id} [get]
func (e ExpenseHandler) GetExpense(c echo.Context) error {
	expenseID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}
	expense, err := e.expenseModel.ReadOne(bson.M{"_id": expenseID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusOK, expense, "expense detail", c)
}

// DeleteExpense godoc
// @Summary Delete a Expense.
// @Description delete expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/expenses/{id} [delete]
func (e ExpenseHandler) DeleteExpense(c echo.Context) error {
	expenseID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	count, err := e.expenseModel.Remove(bson.M{"_id": expenseID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusAccepted, count, "expense removed", c)
}

// UpdateExpense godoc
// @Summary Update an Expense.
// @Description update expense by ID
// @Tags expenses
// @Accept json
// @Produce json
// @Param expense body models.ExpenseInput true "Update Expense"
// @Param id path string true "Expense ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/expenses/{id} [put]
func (e ExpenseHandler) UpdateExpense(c echo.Context) error {
	expenseID, err := objectIDFromStringID(c.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	expInput := new(models.ExpenseInput)

	if err := c.Bind(expInput); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := c.Validate(expInput); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	categoryID, err := objectIDFromStringID(expInput.CategoryID)
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	category, err := e.categoryModel.ReadOne(bson.M{"_id": categoryID})
	if err != nil {
		log.Printf("CATEGORY NOT FOUND ERROR: %v\n", err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}

	userID, err := objectIDFromStringID(expInput.InsertedBy)
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), c)
	}

	user, err := e.userModel.ReadOneUser(bson.M{"_id": userID})
	if err != nil {
		log.Printf("USER NOT FOUND ERROR: %v\n", err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}

	// update fields - title. description, date, category, location, total, status, inserted_by
	update := bson.M{
		"title":       expInput.Title,
		"description": expInput.Description,
		"date":        expInput.Date,
		"category":    category,
		"location":    expInput.Location,
		"total":       expInput.Total,
		"status":      expInput.Status,
		"inserted_by": user,
		"updated_at":  time.Now(),
	}

	count, err := e.expenseModel.UpdateOne(update, bson.M{"_id": expenseID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), c)
	}
	return utils.Data(http.StatusOK, count, "expense updated", c)
}
