package handlers

import (
	"awesome/image-storage-service/service/image-storage/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowPhotos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var photos []entity.Photo
		if err := db.Find(&photos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "view.html", gin.H{
			"photos": photos,
		})
	}
}
