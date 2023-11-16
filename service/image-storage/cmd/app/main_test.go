package main

import (
	"awesome/image-storage-service/service/image-storage/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestConnectToDB_Success(t *testing.T) {
	// Arrange
	cfg := config.Config{
		DB: config.DBConfig{
			Host:     "localhost",
			User:     "testuser",
			Password: "testpass",
			DBName:   "testdb",
			Port:     "5432",
		},
	}

	db := ConnectToDB(cfg)

	assert.NotNil(t, db)
}

func TestGetImageByName_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	db, mock, _ := sqlmock.New()
	gormDB, _ := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})

	mock.ExpectQuery("^SELECT \\* FROM \"photos\"").
		WithArgs("test.jpg").
		WillReturnRows(sqlmock.NewRows([]string{"name", "data"}).
			AddRow("test.jpg", []byte("image data")))

	r.GET("/image/:name", getImageHandler(gormDB))

	req := httptest.NewRequest("GET", "/image/test.jpg", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "image/jpeg", w.Header().Get("Content-Type"))
	assert.Equal(t, []byte("image data"), w.Body.Bytes())
}
