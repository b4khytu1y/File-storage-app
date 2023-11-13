package handlers

import (
	"io/ioutil"
	"net/http"

	"awesome/image-storage-service/service/image-storage/entity" // Путь к вашему пакету entity

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadPhoto(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
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
			Name: file.Filename,
			Data: bytes,
		}

		result := db.Create(&newPhoto)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Фотография успешно загружена", "photo_id": newPhoto.ID})
	}
}
