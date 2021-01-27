package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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

// DeleteProject godoc
// @Summary Delete a Project.
// @Description get project by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/projects/{id} [delete]
func (c ProjectHandler) DeleteProject(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	count, err := c.projectModel.UpdateOne(bson.D{{"is_active", false}}, bson.M{"_id": ID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusAccepted, count, "project removed", e)
}

// GetProjectExpenses godoc
// @Summary Get a Project Details.
// @Description get project details by ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param start query string false "start period with a string representation of date 'YYYY-MM-DD'"
// @Param end query string false "end period with a string representation of date 'YYYY-MM-DD'"
// @Param is_active query string false "is_active to check if the project user is active or not"
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
	filter := models.ProjectDetailsQS{
		Start:    now.BeginningOfMonth(),
		End:      now.EndOfMonth(),
		IsActive: true,
	}
	if len(e.QueryParams()) >= 0 {
		if x, ok := qs["start"]; ok {
			startDate, err := parseDateToFormat("2006-01-02", x[0])
			if err != nil {
				log.Printf("ERROR PARSING STARTDATE: %v\n", err)
				return utils.Error(http.StatusBadRequest, "Specify the Start period", e)
			}
			filter.Start = startDate
		}
		if x, ok := qs["end"]; ok {
			endDate, err := parseDateToFormat("2006-01-02", x[0])
			if err != nil {
				log.Printf("ERROR PARSING ENDDATE: %v\n", err)
				return utils.Error(http.StatusBadRequest, "Specify the End period", e)
			}

			filter.End = endDate
		}

		if x, ok := qs["is_active"]; ok {
			b, err := strconv.ParseBool(x[0])
			if err != nil {
				log.Printf("INVALID QUERY PARAM PASSED: %v\n", err)
				return utils.Error(http.StatusBadRequest, err.Error(), e)
			}
			filter.IsActive = b
		}
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

	b, err := ioutil.ReadAll(e.Request().Body)
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	p := new(models.ProjectUser)
	err = json.Unmarshal(b, &p)
	if err != nil {
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
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

// GetProjectUsers godoc
// @Summary Create a Project User.
// @Description create a project user.
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 500 {object} utils.Response
// @Router /api/v1/projects/{id}/users [get]
func (c ProjectHandler) GetProjectUsers(e echo.Context) error {
	ID, err := objectIDFromStringID(e.Param("id"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	qs := e.QueryParams()
	filter := bson.D{
		{"project_id", ID},
	}
	if len(e.QueryParams()) >= 0 {
		if x, ok := qs["is_active"]; ok {
			b, err := strconv.ParseBool(x[0])
			if err != nil {
				log.Printf("INVALID QUERY PARAM PASSED: %v\n", err)
				return utils.Error(http.StatusBadRequest, err.Error(), e)
			}
			filter = bson.D{
				{"project_id", ID},
				{"is_active", b},
			}
		}
	}

	user, err := c.projectModel.ReadAllProjectUser(filter)
	if err != nil {
		log.Printf("RESPONSE ERROR: %v\n", err)
		return utils.Error(http.StatusInternalServerError, err.Error(), e)
	}

	return utils.Data(http.StatusOK, user, "project user details", e)
}

// GetProjectUser godoc
// @Summary Get a Project User.
// @Description get project by ID and project_user ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param userId path string true "Project User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/projects/{id}/users/{userId} [get]
func (c ProjectHandler) GetProjectUser(e echo.Context) error {
	userID, err := objectIDFromStringID(e.Param("userId"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	count, err := c.projectModel.ReadOneProjectUser(bson.M{"_id": userID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusAccepted, count, "project user removed", e)
}

// DeleteProjectUser godoc
// @Summary Delete a Project User.
// @Description delete project by ID and project_user ID
// @Tags projects
// @Accept json
// @Produce json
// @Param id path string true "Project ID"
// @Param userId path string true "Project User ID"
// @Success 200 {object} utils.Response
// @Failure 400 {object} utils.Response
// @Failure 404 {object} utils.Response
// @Router /api/v1/projects/{id}/users/{userId} [delete]
func (c ProjectHandler) DeleteProjectUser(e echo.Context) error {
	userID, err := objectIDFromStringID(e.Param("userId"))
	if err != nil {
		return utils.Error(http.StatusBadRequest, err.Error(), e)
	}

	// soft delete
	count, err := c.projectModel.UpdateOneProjectUser(bson.M{"is_active": false}, bson.M{"_id": userID})
	if err != nil {
		log.Println(err)
		return utils.Error(http.StatusNotFound, err.Error(), e)
	}
	return utils.Data(http.StatusAccepted, count, "project user removed", e)
}
