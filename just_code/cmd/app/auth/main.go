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

func main() {

	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}

	//Database
	db := config.ConnectionDB(&loadConfig)
	validate := validator.New()
	db.Table("users").AutoMigrate(&model.Users{})

	// log.Println("Starting auto migration files")
	// err = db.Table("files").AutoMigrate(&model.Users{}, &model.FileModel{})
	// if err != nil {
	// 	log.Fatalf("Auto migration failed: %v", err)
	// }

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
