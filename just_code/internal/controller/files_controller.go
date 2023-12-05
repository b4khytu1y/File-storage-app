package controller

import (
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/service"
	"golang-jwttoken/pkg/utils"
	"io/ioutil" //nolint:staticcheck
	"log"
	"strconv"

	_ "golang-jwttoken/docs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type FileController struct {
	fileService service.FileService
}

type FileControllerBuilder struct {
	fileService service.FileService
}

func NewFileControllerBuilder() *FileControllerBuilder {
	return &FileControllerBuilder{}
}

func (b *FileControllerBuilder) SetFileService(fileService service.FileService) *FileControllerBuilder {
	b.fileService = fileService
	return b
}

func (b *FileControllerBuilder) Build() *FileController {
	return &FileController{
		fileService: b.fileService,
	}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Uploads a new file to the server
// @Tags files
// @Accept  multipart/form-data
// @Produce  json
// @Param file formData file true "File to upload"
// @Security ApiKeyAuth
// @Success 200 {object} gin.H "Файл успешно загружен"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /files [post]
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
// @Description Retrieves a file by its ID
// @Tags files
// @Produce  json
// @Param id path int true "File ID"
// @Security ApiKeyAuth
// @Success 200 {object} gin.H "File data"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
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
// @Summary Retrieve files of a user
// @Description Retrieves all files uploaded by the authenticated user
// @Tags files
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {array} model.FileModel
// @Failure 401 "Unauthorized"
// @Failure 500 "Internal Server Error"
// @Router /user/files [get]
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
// @Summary Update file information
// @Description Updates the metadata of a specified file
// @Tags files
// @Accept  json
// @Produce  json
// @Param id path int true "File ID"
// @Security ApiKeyAuth
// @Success 200 "Файл успешно обновлен"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
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
// @Description Deletes a file by its ID
// @Tags files
// @Produce  json
// @Param id path int true "File ID"
// @Security ApiKeyAuth
// @Success 200 "Файл успешно удален"
// @Failure 400 "Bad Request"
// @Failure 401 "Unauthorized"
// @Failure 403 "Forbidden"
// @Failure 404 "Not Found"
// @Failure 500 "Internal Server Error"
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
