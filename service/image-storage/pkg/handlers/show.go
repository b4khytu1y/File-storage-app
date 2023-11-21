package handlers

import (
	"awesome/image-storage-service/service/image-storage/internal/entity"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func ShowPhotos(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не аутентифицирован"})
			return
		}

		var photos []entity.Photo
		if err := db.Where("user_id = ?", userID).Find(&photos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "view.html", gin.H{
			"photos": photos,
		})
	}
}
