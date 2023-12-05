package main

import (
	"context"
	"golang-jwttoken/config"
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
	"golang-jwttoken/internal/router"
	"golang-jwttoken/internal/service"
	"os"
	"os/signal"
	"syscall"

	"golang-jwttoken/internal/controller"

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
		log.Fatal("Could not load environment variables", err)
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
	fileController := controller.NewFileControllerBuilder().
		SetFileService(fileService).
		Build()
	routes := router.NewRouter(userRepository, authenticationController, usersController, fileController)

	server := &http.Server{
		Addr:           ":" + loadConfig.ServerPort,
		Handler:        routes,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe(): %s", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
