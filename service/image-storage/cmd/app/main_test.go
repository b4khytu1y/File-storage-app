package main

// import (
// 	"awesome/image-storage-service/service/image-storage/config"
// 	"awesome/image-storage-service/service/image-storage/entity"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/DATA-DOG/go-sqlmock"
// 	"github.com/gin-gonic/gin"
// 	"github.com/stretchr/testify/assert"
// 	"gorm.io/driver/postgres"
// 	"gorm.io/gorm"
// )

// func getImageHandler(db *gorm.DB) gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		var photo entity.Photo
// 		name := c.Param("name")

// 		// Retrieve the photo by name from the database
// 		result := db.Where("name = ?", name).First(&photo)
// 		if result.Error != nil {
// 			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
// 				c.JSON(http.StatusNotFound, gin.H{"error": "Image not found"})
// 				return
// 			}
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
// 			return
// 		}

// 		// Respond with the photo data
// 		c.Data(http.StatusOK, "image/jpeg", photo.Data)
// 	}
// }
// func TestGetImageByName_Success(t *testing.T) {
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()

// 	db, mock, _ := sqlmock.New()
// 	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

// 	mock.ExpectQuery("^SELECT \\* FROM \"photos\"").
// 		WithArgs("test.jpg").
// 		WillReturnRows(sqlmock.NewRows([]string{"name", "data"}).
// 			AddRow("test.jpg", []byte("image data")))

// 	r.GET("/image/:name", getImageHandler(gormDB))

// 	req := httptest.NewRequest("GET", "/image/test.jpg", nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusOK, w.Code)
// 	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
// 	assert.Equal(t, []byte("image data"), w.Body.Bytes())
// }

// func TestGetImageByName_NotFound(t *testing.T) {
// 	// Arrange
// 	gin.SetMode(gin.TestMode)
// 	r := gin.Default()

// 	db, mock, _ := sqlmock.New()
// 	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

// 	mock.ExpectQuery("^SELECT \\* FROM \"photos\"").
// 		WithArgs("nonexistent.jpg").
// 		WillReturnError(gorm.ErrRecordNotFound)

// 	r.GET("/image/:name", getImageHandler(gormDB))

// 	req := httptest.NewRequest("GET", "/image/nonexistent.jpg", nil)
// 	w := httptest.NewRecorder()
// 	r.ServeHTTP(w, req)

// 	assert.Equal(t, http.StatusNotFound, w.Code)
// }

// func TestGetImageByName_DBError(t *testing.T) {

// 	mock.ExpectQuery("^SELECT \\* FROM \"photos\"").
// 		WithArgs("test.jpg").
// 		WillReturnError(errors.New("some db error"))

// }
// func TestConnectToDB_Failure(t *testing.T) {
// 	cfg := config.Config{
// 		DB: config.DBConfig{
// 			Host:     "invalid_host",
// 			User:     "invalid_user",
// 			Password: "invalid_password",
// 			DBName:   "invalid_dbname",
// 			Port:     "invalid_port",
// 		},
// 	}

// 	assert.PanicsWithValue(t, "failed to connect database", func() { ConnectToDB(cfg) }, "The code did not panic as expected")
// }
