package controller

import (
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/service"
	"golang-jwttoken/pkg/utils"
	"io/ioutil"
	"log"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type FileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) *FileController {
	return &FileController{fileService: fileService}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a new file
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param file formData file true "Upload file"
// @Success 200 {object} map[string]interface{}
// @Router /files/upload [post]
func (fc *FileController) UploadFile(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := utils.ValidateToken(tokenString, "Secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "http: no such file"})
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		log.Printf("Ошибка при открытии файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при открытии файла"})
		return
	}
	defer openedFile.Close()

	fileContents, err := ioutil.ReadAll(openedFile)
	if err != nil {
		log.Printf("Ошибка при чтении файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при чтении файла"})
		return
	}

	fileModel := &model.FileModel{
		UserID:      userID,
		Name:        file.Filename,
		ContentType: file.Header.Get("Content-Type"),
		Size:        file.Size,
		Content:     fileContents,
	}

	err = fc.fileService.SaveFile(fileModel)
	if err != nil {
		log.Printf("Ошибка при сохранении файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении файла"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно загружен"})
}

// GetFile godoc
// @Summary Retrieve a file
// @Description Get details of a file by file ID
// @Tags files
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path int true "File ID"
// @Success 200 {object} model.FileModel
// @Router /files/{id} [get]
func (fc *FileController) GetFile(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := utils.ValidateToken(tokenString, "Secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}

	fileIDString := c.Param("id")
	if fileIDString == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID файла отсутствует"})
		return
	}

	fileID, err := strconv.Atoi(fileIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID файла"})
		return
	}

	file, err := fc.fileService.GetFileByID(fileID)
	if err != nil {
		log.Printf("Ошибка при получении файла: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Файл не найден"})
		return
	}

	user, err := fc.fileService.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка при получении данных пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}

	if user.IsAdmin != "1" && file.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к этому файлу"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"file": file})
}

// GetUserFiles godoc
// @Summary Retrieve all files for a user
// @Description Get all files uploaded by a user
// @Tags files
// @Produce json
// @Param Authorization header string true "Authorization"
// @Success 200 {array} model.FileModel
// @Router /files/user [get]
func (fc *FileController) GetUserFiles(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := utils.ValidateToken(tokenString, "Secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}

	user, err := fc.fileService.GetUserByID(userID)
	if err != nil {
		log.Printf("Ошибка при получении данных пользователя: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}

	var isAdminString string
	if user.IsAdmin == "1" {
		isAdminString = "1"
	} else {
		isAdminString = "0"
	}

	userFiles, err := fc.fileService.GetFilesByUserID(userID, isAdminString)
	if err != nil {
		log.Printf("Ошибка при получении файлов: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении файлов"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": userFiles})
}

// UpdateFile godoc
// @Summary Update a file
// @Description Update file details
// @Tags files
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path int true "File ID"
// @Success 200 {object} map[string]interface{}
// @Router /files/{id} [put]
func (fc *FileController) UpdateFile(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := utils.ValidateToken(tokenString, "Secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}

	fileIDString := c.Param("id")
	fileID, err := strconv.Atoi(fileIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID файла"})
		return
	}

	existingFile, err := fc.fileService.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден"})
		return
	}

	user, err := fc.fileService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}

	if user.IsAdmin != "1" && existingFile.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к обновлению файла"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно обновлен"})
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file by file ID
// @Tags files
// @Produce json
// @Param Authorization header string true "Authorization"
// @Param id path int true "File ID"
// @Success 200 {object} map[string]interface{}
// @Router /files/{id} [delete]
func (fc *FileController) DeleteFile(c *gin.Context) {
	tokenString := c.GetHeader("Authorization")
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	userID, err := utils.ValidateToken(tokenString, "Secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Невалидный токен"})
		return
	}

	fileIDString := c.Param("id")
	fileID, err := strconv.Atoi(fileIDString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат ID файла"})
		return
	}

	existingFile, err := fc.fileService.GetFileByID(fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Файл не найден"})
		return
	}

	user, err := fc.fileService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при получении данных пользователя"})
		return
	}

	if user.IsAdmin != "1" && existingFile.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Нет доступа к удалению файла"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Файл успешно удален"})
}
func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return r
}
