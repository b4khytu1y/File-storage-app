package controller

import (
	controller "awesome/image-storage-service/service/image-storage/internal/controller/handler"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitializeRoutes(router *gin.Engine, db *gorm.DB, secretKey string) {

	uploadImageHandler := controller.UploadImageHandler(db)
	getImageHandler := controller.GetImageHandler(db)
	viewImagesHandler := controller.ViewImagesHandler(db)

	protectedRoutes := router.Group("/")
	{
		protectedRoutes.POST("/upload", uploadImageHandler)
		protectedRoutes.GET("/image/:name", getImageHandler)
		protectedRoutes.GET("/view", viewImagesHandler)
	}
}
