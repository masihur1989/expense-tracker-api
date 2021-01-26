package handler

import (
	"log"
	"net/http"
	"time"

	"github.com/jinzhu/now"
	"github.com/labstack/echo/v4"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	"github.com/masihur1989/expense-tracker-api/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ProjectHandler godoc
type ProjectHandler struct {
	projectModel models.ProjectModeler
}

// NewProjectHandler godoc
func NewProjectHandler(pm models.ProjectModeler) ProjectHandler {
	return ProjectHandler{pm}
}

// CreateProject godoc
// @Summary Create Project.
// @Description create project.
// @Tags projects
// @Accept json
// @Produce json
// @Param project body models.Project true "Create Project"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects [post]
func (c ProjectHandler) CreateProject(e echo.Context) error {
	p := new(models.Project)
	if err := e.Bind(p); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := e.Validate(p); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	// fill the nil values
	p.ID = primitive.NewObjectID()
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	id, err := c.projectModel.Insert(p)

	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}

	return utils.Data(http.StatusCreated, id, "projects created", e)
}

// GetProjects godoc
// get all the categories
// QueryParams accepted are name
// and only single value is used. TODO needs to find a way to accept array of query pamra values
// @Summary Get Projects.
// @Description get projects
// @Tags projects
// @Accept json
// @Produce json
// @Param name query string false "name search by name"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects [get]
func (c ProjectHandler) GetProjects(e echo.Context) error {
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

	pats, err := c.projectModel.ReadAll(filter)
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}
	return utils.Data(http.StatusOK, pats, "project details", e)
}

// GetProject godoc
// @Summary Get an Project.
// @Description get project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects/{id} [get]
func (c ProjectHandler) GetProject(e echo.Context) error {
	projectID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}
	expense, err := c.projectModel.ReadOne(bson.M{"_id": projectID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusOK, expense, "project detail", e)
}

// GetProjectExpenses godoc
// @Summary Get a Project Expenses with Users.
// @Description get project details by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param start query string false "start period with a string representation of date 'YYYY-MM-DD'"
// @Param end query string false "end period with a string representation of date 'YYYY-MM-DD'"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects/{id}/details [get]
func (c ProjectHandler) GetProjectExpenses(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}
	qs := e.QueryParams()
	filter := models.ProjectExpenseQS{}
	if len(e.QueryParams()) == 0 {
		filter.Start = now.BeginningOfMonth()
		filter.End = now.EndOfMonth()
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
			return utils.Error(http.StatusBadRequest, "Specify the Start period", e)
		}

		endDate, err = parseDateToFormat("2006-01-02", f["end"])
		if err != nil {
			log.Printf("ERROR PARSING ENDDATE: %v\n", err)
			return utils.Error(http.StatusBadRequest, "Specify the End period", e)
		}

		filter.Start = startDate
		filter.End = endDate
	}

	pats, err := c.projectModel.LookupProjectDetails(bson.D{{"_id", ID}}, filter)
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}
	return utils.Data(http.StatusOK, pats, "complete project details", e)
}

// CreateProjectUser godoc
// @Summary Create a Project User.
// @Description create a project user.
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param projectUser body models.ProjectUser true "Add Project User"
// @Success 201 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects/{id}/users [post]
func (c ProjectHandler) CreateProjectUser(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}
	p := new(models.ProjectUser)
	if err := e.Bind(p); err != nil {
		log.Printf("ECHO BINDING ERROR: %v\n", err)
		return err
	}

	if err := e.Validate(p); err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	p.ID = primitive.NewObjectID()
	p.ProjectID = ID
	p.CreatedAt = time.Now()
	p.UpdatedAt = time.Now()

	id, err := c.projectModel.InsertProjectUser(p)
	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}

	return utils.Data(http.StatusCreated, id, "project user created", e)
}
