package service

import (
	"errors"
	"golang-jwttoken/config"
	"golang-jwttoken/internal/data/request"
	"golang-jwttoken/internal/helper"
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
	"golang-jwttoken/pkg/utils"

	"github.com/go-playground/validator/v10"
)

type AuthenticationServiceImpl struct {
	UsersRepository repository.UsersRepository
	Validate        *validator.Validate
}

func NewAuthenticationServiceImpl(usersRepository repository.UsersRepository, validate *validator.Validate) AuthenticationService {
	return &AuthenticationServiceImpl{
		UsersRepository: usersRepository,
		Validate:        validate,
	}
}

func (a *AuthenticationServiceImpl) Login(users request.LoginRequest) (string, error) {
	new_users, users_err := a.UsersRepository.FindByUsername(users.Username)
	if users_err != nil {
		return "", errors.New("invalid username or Password")
	}

	config, _ := config.LoadConfig(".")

	verify_error := utils.VerifyPassword(new_users.Password, users.Password)
	if verify_error != nil {
		return "", errors.New("invalid username or Password")
	}

	token, errToken := utils.GenerateToken(new_users.Id, config.TokenExpiresIn, config.TokenSecret)
	if errToken != nil {
		return "", errToken
	}
	return token, nil
}

func (a *AuthenticationServiceImpl) Register(users request.CreateUsersRequest) {

	hashedPassword, err := utils.HashPassword(users.Password)
	helper.ErrorPanic(err)

	newUser := model.Users{
		Username: users.Username,
		Email:    users.Email,
		Password: hashedPassword,
		IsAdmin:  users.IsAdmin,
	}
	a.UsersRepository.Save(newUser)
}
