package main

import (
	"golang-jwttoken/config"
	"golang-jwttoken/internal/controller"
	"golang-jwttoken/internal/helper"
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
	"golang-jwttoken/internal/router"
	"golang-jwttoken/internal/service"

	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

// @title           Swagger API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	db := config.ConnectionDB(&loadConfig)
	validate := validator.New()
	db.Table("users").AutoMigrate(&model.Users{}) //nolint:all

	userRepository := repository.NewUsersRepositoryImpl(db)
	fileRepository := repository.NewFileRepositoryImpl(db)

	fileService := service.NewFileServiceImpl(fileRepository, userRepository)

	authenticationService := service.NewAuthenticationServiceImpl(userRepository, validate)

	authenticationController := controller.NewAuthenticationController(authenticationService)
	usersController := controller.NewUsersController(userRepository)
	fileController := controller.NewFileController(fileService)
	routes := router.NewRouter(userRepository, authenticationController, usersController, fileController)

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	server_err := server.ListenAndServe()
	helper.ErrorPanic(server_err)

}
