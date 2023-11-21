package handlers

import (
	"awesome/image-storage-service/service/image-storage/internal/entity"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadPhoto(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
			return
		}
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer src.Close()

		bytes, err := ioutil.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newPhoto := entity.Photo{
			UserID: userID.(int),
			Name:   file.Filename,
			Data:   bytes,
		}

		if result := db.Create(&newPhoto); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Фотография успешно загружена"})
	}
}
func getImageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var photo entity.Photo
		name := c.Param("name")

		// Retrieve the photo by name from the database
		result := db.Where("name = ?", name).First(&photo)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		// Respond with the photo data
		c.Data(http.StatusOK, "image/jpeg", photo.Data)
	}
}
