package controller

import (
	"fmt"

	"golang-jwttoken/internal/data/request"
	"golang-jwttoken/internal/data/response"
	"golang-jwttoken/internal/helper"
	"golang-jwttoken/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthenticationController struct {
	authenticationService service.AuthenticationService
}

func NewAuthenticationController(service service.AuthenticationService) *AuthenticationController {
	return &AuthenticationController{authenticationService: service}
}

// Login godoc
// @Summary Login user
// @Description Authenticate user and return token
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param loginRequest body request.LoginRequest true "Login Information"
// @Success 200 {object} response.Response{Data=response.LoginResponse}
// @Failure 400 {object} response.Response "Invalid username or password"
// @Router /login [post]
func (controller *AuthenticationController) Login(ctx *gin.Context) {
	loginRequest := request.LoginRequest{}
	err := ctx.ShouldBindJSON(&loginRequest)
	helper.ErrorPanic(err)

	token, err_token := controller.authenticationService.Login(loginRequest)
	fmt.Println(err_token)
	if err_token != nil {
		webResponse := response.Response{
			Code:    http.StatusBadRequest,
			Status:  "Bad Request",
			Message: "Invalid username or password",
		}
		ctx.JSON(http.StatusBadRequest, webResponse)
		return
	}

	resp := response.LoginResponse{
		TokenType: "Bearer",
		Token:     token,
	}

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully log in!",
		Data:    resp,
	}

	// ctx.SetCookie("token", token, config.TokenMaxAge*60, "/", "localhost", false, true)
	ctx.JSON(http.StatusOK, webResponse)
}

// Register godoc
// @Summary Register user
// @Description Register a new user
// @Tags authentication
// @Accept  json
// @Produce  json
// @Param createUsersRequest body request.CreateUsersRequest true "User Registration Information"
// @Success 200 {object} response.Response "Successfully created user!"
// @Router /register [post]
func (controller *AuthenticationController) Register(ctx *gin.Context) {
	createUsersRequest := request.CreateUsersRequest{}
	err := ctx.ShouldBindJSON(&createUsersRequest)
	helper.ErrorPanic(err)

	controller.authenticationService.Register(createUsersRequest)

	webResponse := response.Response{
		Code:    200,
		Status:  "Ok",
		Message: "Successfully created user!",
		Data:    nil,
	}

	ctx.JSON(http.StatusOK, webResponse)
}
