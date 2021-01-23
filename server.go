package main

import (
	"log"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	_ "github.com/joho/godotenv/autoload"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/masihur1989/expense-tracker-api/docs" // you need to update github.com/rizalgowandy/go-swag-sample with your own project path
	db "github.com/masihur1989/expense-tracker-api/internal/db"
	"github.com/masihur1989/expense-tracker-api/internal/handler"
	"github.com/masihur1989/expense-tracker-api/internal/models"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gopkg.in/go-playground/validator.v9"
	en_translations "gopkg.in/go-playground/validator.v9/translations/en"

	customMiddleware "github.com/masihur1989/expense-tracker-api/internal/middleware"
)

// @title Expense Tracker Swagger API
// @version 1.0
// @description This is a sample server server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:1323
// @BasePath /
// @schemes http https
func main() {
	// setup echo
	e := SetupEcho()
	// routes
	e.GET("/", handler.Ping)
	e.GET("/docs/*", echoSwagger.WrapHandler)

	// initialize db client
	client, err := db.GetClient()
	if err != nil {
		log.Panicf("DB CONNECTION ERROR: %f", err)
	}
	// models
	userModel := models.NewUserModelImpl(client)
	categoryModel := models.NewCategoryModel(client)
	// route versioning /api/v1
	g := e.Group("/api/v1")
	// handlers
	userHandler := handler.NewUserHandler(userModel)
	categoryHandler := handler.NewCategoryHandler(categoryModel)
	// users routes
	g.GET("/users/", userHandler.GetUsers)
	g.GET("/users/:id", userHandler.GetUser)
	g.POST("/users", userHandler.CreateUser)
	g.PUT("/users/:id", userHandler.UpdateUser)
	g.DELETE("/users/:id", userHandler.DeleteUser)
	// categories routes
	g.GET("/categories", categoryHandler.GetCategories)
	g.POST("/categories", categoryHandler.CreateCategory)
	g.PUT("/categories/:id", categoryHandler.UpdateCategory)
	g.DELETE("/categories/:id", categoryHandler.DeleteCategory)

	e.Logger.Fatal(e.Start(":1323"))
}

/*******************************************************************************************************
										SETUP FUNCTIONS
*******************************************************************************************************/

// SetupEcho set echo server
func SetupEcho() *echo.Echo {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// custom middlewares
	e.Use(customMiddleware.RequestHeaders())
	// setup validator
	trans := SetupTranslator()
	validate := SetupCustomValidator(trans)
	e.Validator = &models.Validator{Validator: validate, Trans: trans}
	return e
}

// SetupTranslator set the translator
func SetupTranslator() ut.Translator {
	translator := en.New()
	uni := ut.New(translator, translator)

	trans, found := uni.GetTranslator("en")
	log.Printf("SetupTranslator: trans %v\n", trans)
	if !found {
		log.Fatalln("translator not found")
	}
	return trans
}

// SetupCustomValidator set the custom validator
func SetupCustomValidator(trans ut.Translator) *validator.Validate {
	v := validator.New()

	if err := en_translations.RegisterDefaultTranslations(v, trans); err != nil {
		log.Fatal(err)
	}

	return v
}
