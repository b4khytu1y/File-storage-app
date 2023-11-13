package main

import (
	"awesome/image-storage-service/service/image-storage/config"
	"awesome/image-storage-service/service/image-storage/entity"
	"fmt"
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
	viper.SetConfigFile("D:\\go\\image-storage-service\\image-storage-service\\service\\image-storage\\config\\config.yaml")
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

	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filePath := "uploads/" + file.Filename

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		newPhoto := entity.Photo{}
		if result := db.Create(&newPhoto); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		c.String(http.StatusOK, "Файл успешно загружен")
	})

	router.LoadHTMLGlob("../../templates/*")

	router.Run(":8080")
}
