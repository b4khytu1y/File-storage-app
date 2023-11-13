package main

import (
	"awesome/image-storage-service/service/image-storage/entity"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectToDB(cfg Config) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&entity.Photo{})
	return db
}

func main() {
	cfg, err := LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	db := ConnectToDB(*cfg)

	// Настройка Gin
	router := gin.Default()

	// Статические файлы
	router.Static("/assets", "./assets")

	// Маршрут для отображения формы загрузки
	router.GET("/upload", func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	// Обработчик загрузки изображений
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Здесь можно добавить логику для сохранения файла с использованием GORM
		// ...

		c.String(http.StatusOK, "Файл успешно загружен")
	})

	// Загрузка HTML-шаблонов
	router.LoadHTMLGlob("templates/*")

	// Запуск сервера
	router.Run(":8080")
}
