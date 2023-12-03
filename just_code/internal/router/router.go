package router

import (
	"golang-jwttoken/internal/controller"
	"golang-jwttoken/internal/middleware"
	"golang-jwttoken/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter(userRepository repository.UsersRepository, authenticationController *controller.AuthenticationController, usersController *controller.UserController, fileController *controller.FileController) *gin.Engine {
	service := gin.Default()

	service.GET("", func(context *gin.Context) {
		context.JSON(http.StatusOK, "welcome home")
	})

	service.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	router := service.Group("/api")
	authenticationRouter := router.Group("/authentication")
	authenticationRouter.POST("/register", authenticationController.Register)
	authenticationRouter.POST("/login", authenticationController.Login)

	usersRouter := router.Group("/users")
	usersRouter.GET("", middleware.DeserializeUser(userRepository), usersController.GetUsers)
	usersRouter.GET("/:id", middleware.DeserializeUser(userRepository), usersController.GetUser)
	usersRouter.PUT("/:id", middleware.DeserializeUser(userRepository), usersController.UpdateUser)
	usersRouter.DELETE("/:id", middleware.DeserializeUser(userRepository), usersController.DeleteUser)

	filesRouter := router.Group("/files")
	filesRouter.GET("", middleware.DeserializeUser(userRepository), fileController.GetUserFiles)
	filesRouter.GET("/:id", middleware.DeserializeUser(userRepository), fileController.GetFile)
	filesRouter.PUT("/:id", middleware.DeserializeUser(userRepository), fileController.GetUserFiles)
	filesRouter.DELETE("/:id", middleware.DeserializeUser(userRepository), fileController.GetUserFiles)
	filesRouter.POST("/upload", middleware.DeserializeUser(userRepository), fileController.UploadFile)

	return service
}
