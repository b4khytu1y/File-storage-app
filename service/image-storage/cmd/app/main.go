package main

import (
	"awesome/image-storage-service/service/image-storage/config"
	"awesome/image-storage-service/service/image-storage/entity"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Photo{})
	return db
}

func main() {

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурации: %s", err)
	}

	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Ошибка при десериализации конфигурации: %s", err)
	}

	db := ConnectToDB(cfg)
	router := gin.Default()

	router.Static("/assets", "./assets")

	router.GET("/image/:name", func(c *gin.Context) {
		name := c.Param("name")

		var photo entity.Photo
		if result := db.Where("name = ?", name).First(&photo); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Изображение не найдено"})
			} else {
				c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			}
			return
		}

		c.Writer.Header().Set("Content-Type", "image/jpeg")
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write(photo.Data)
	})
	router.GET("/view", func(c *gin.Context) {
		var photos []entity.Photo
		if err := db.Find(&photos).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.HTML(http.StatusOK, "view.html", gin.H{"photos": photos})
	})

	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	router.POST("/upload", func(c *gin.Context) {
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

		data, err := ioutil.ReadAll(src)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newPhoto := entity.Photo{
			Name: file.Filename,
			Data: data,
		}

		if result := db.Create(&newPhoto); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.String(http.StatusOK, "Файл успешно загружен")
	})

	router.LoadHTMLGlob("../../templates/*")

	router.Run(":8080")
}
