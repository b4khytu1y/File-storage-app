package controller

import (
	"awesome/image-storage-service/service/image-storage/internal/entity"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func UploadImageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить файл"})
			return
		}

		src, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось открыть файл"})
			return
		}
		defer src.Close()

		data, err := ioutil.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось прочитать файл"})
			return
		}

		userID, _ := c.Get("userID")

		newPhoto := entity.Photo{
			UserID: userID.(int),
			Name:   file.Filename,
			Data:   data,
		}

		if err := db.Create(&newPhoto).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить изображение"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Изображение успешно загружено"})
	}
}

func GetImageHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")

		var photo entity.Photo
		if err := db.Where("name = ?", name).First(&photo).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Изображение не найдено"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при поиске изображения"})
			return
		}

		c.Header("Content-Type", "image/jpeg")
		c.Header("Content-Length", fmt.Sprintf("%d", len(photo.Data)))
		c.Writer.Write(photo.Data)
	}
}

func ViewImagesHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, _ := c.Get("userID")

		var photos []entity.Photo
		if err := db.Where("user_id = ?", userID).Find(&photos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении списка изображений"})
			return
		}

		c.HTML(http.StatusOK, "view.html", gin.H{"photos": photos})
	}
}
