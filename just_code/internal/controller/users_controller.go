package controller

import (
	_ "golang-jwttoken/docs"
	"golang-jwttoken/internal/data/response"
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userRepository repository.UsersRepository
}

func NewUsersController(repository repository.UsersRepository) *UserController {
	return &UserController{userRepository: repository}
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieves a list of users, admin can retrieve all, others can retrieve only their information
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response "Successfully fetch user data"
// @Router /users [get]
func (controller *UserController) GetUsers(ctx *gin.Context) {
	currentUser := ctx.GetString("currentUser")
	isAdmin := ctx.GetString("isAdmin")

	var users []model.Users

	if isAdmin == "1" {
		users = controller.userRepository.FindAll()
	} else {
		user, err := controller.userRepository.FindByUsername(currentUser)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Пользователь не найден"})
			return
		}
		users = []model.Users{user}
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully fetch user data!",
		Data:    users,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// GetUser godoc
// @Summary Get a single user
// @Description Retrieves details of a specific user, accessible to admin or the user themselves
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response "Successfully fetch user data"
// @Failure 404 "User not found"
// @Failure 400 "Invalid user ID format"
// @Router /users/{id} [get]
func (controller *UserController) GetUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Неверный формат идентификатора пользователя"})
		return
	}
	currentUser := ctx.GetString("currentUser")
	isAdmin := ctx.GetString("isAdmin")

	var users []model.Users

	if isAdmin == "1" || currentUser == userIDStr {
		user, err := controller.userRepository.FindById(userID)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": "Пользователь не найден"})
			return
		}
		users = []model.Users{user}
	} else {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "У вас нет прав на просмотр этого пользователя"})
		return
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully fetch user data!",
		Data:    users,
	}

	ctx.JSON(http.StatusOK, webResponse)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update user details, accessible to admin or the user themselves
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Param updatedUser body model.Users true "Updated User Information"
// @Success 200 "User successfully updated"
// @Failure 400 "Invalid user ID format or data format"
// @Failure 403 "No rights to update this user"
// @Router /users/{id} [put]
func (controller *UserController) UpdateUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Неверный формат идентификатора пользователя"})
		return
	}
	currentUser := ctx.GetString("currentUser")
	isAdmin := ctx.GetString("isAdmin")

	if isAdmin != "1" && currentUser != userIDStr {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "У вас нет прав на обновление этого пользователя"})
		return
	}

	var updatedUser model.Users
	if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Неверный формат данных"})
		return
	}

	if err := controller.userRepository.UpdateUser(userID, &updatedUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Не удалось обновить пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Пользователь успешно обновлен"})
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a specific user, accessible to admin or the user themselves
// @Tags users
// @Accept  json
// @Produce  json
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 "User successfully deleted"
// @Failure 400 "Invalid user ID format"
// @Failure 403 "No rights to delete this user"
// @Router /users/{id} [delete]
func (controller *UserController) DeleteUser(ctx *gin.Context) {
	userIDStr := ctx.Param("id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Неверный формат идентификатора пользователя"})
		return
	}
	currentUser := ctx.GetString("currentUser")
	isAdmin := ctx.GetString("isAdmin")

	if isAdmin != "1" && currentUser != userIDStr {
		ctx.JSON(http.StatusForbidden, gin.H{"status": "fail", "message": "У вас нет прав на удаление этого пользователя"})
		return
	}

	if err := controller.userRepository.DeleteUser(userID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Не удалось удалить пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Пользователь успешно удален"})
}
