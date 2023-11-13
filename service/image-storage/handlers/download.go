package handlers

import (
	"awesome/image-storage-service/service/image-storage/entity"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func DownloadPhoto(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный ID"})
			return
		}

		var photo entity.Photo
		err = db.QueryRow("SELECT id, name, data FROM photos WHERE id = $1", id).Scan(&photo.ID, &photo.Name, &photo.Data)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Фотография не найдена"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			return
		}

		c.Writer.Header().Set("Content-Type", "application/octet-stream")
		c.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", photo.Name))
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(photo.Data)
	}
}
