package main

import (
	"awesome/image-storage-service/service/image-storage/config"
	"awesome/image-storage-service/service/image-storage/internal/entity"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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
	router.Use(TokenAuthMiddleware())

	router.Static("/assets", "./assets")

	router.GET("/image/:name", TokenAuthMiddleware(), func(c *gin.Context) {
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

	router.GET("/upload", TokenAuthMiddleware(), func(c *gin.Context) {
		c.HTML(http.StatusOK, "upload.html", nil)
	})

	router.POST("/upload", TokenAuthMiddleware(), func(c *gin.Context) {
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

		userID := c.MustGet("userID").(int)
		newPhoto := entity.Photo{
			UserID: userID,
			Name:   file.Filename,
			Data:   data,
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

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		const Bearer_schema = "Bearer "
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "API token is required"})
			return
		}
		tokenString := header[len(Bearer_schema):]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte("your_secret_key"), nil
		})

		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			c.Set("userID", claims["user_id"])
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
	}
}

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
