package main

import (
	"log"
	"net/http"

	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/masihur1989/expense-tracker-api/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	db "github.com/masihur1989/expense-tracker-api/internal/db"
	"github.com/masihur1989/expense-tracker-api/internal/handler"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopkg.in/go-playground/validator.v9"
)

// @title Palki-CMS Swagger API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
// @schemes http https
func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	validate := validator.New()
	e.Validator = &models.Validator{Validator: validate}

	client, err := db.GetClient()
	if err != nil {
		log.Panicf("DB CONNECTION ERROR: %f", err)
	}

	// routes
	e.GET("/", Ping)
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// route versioning /api/v1
	g := e.Group("/api/v1")
	handler.NewUserHandler(g, client)

	e.Logger.Fatal(e.Start(":1323"))
}

// Ping godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Server is up and running",
	})
}

// DbPing check for db to be connected
func DbPing(c echo.Context) error {
	_, err := db.GetClient()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"message": "DB Connection Error",
		})
	}
	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "DB Connected",
	})
}
