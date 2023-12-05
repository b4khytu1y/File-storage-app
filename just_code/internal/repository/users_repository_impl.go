package repository

import (
	"errors"

	"golang-jwttoken/internal/data/request"
	"golang-jwttoken/internal/helper"
	"golang-jwttoken/internal/model"

	"gorm.io/gorm"
)

type UsersRepositoryImpl struct {
	Db *gorm.DB
}

func NewUsersRepositoryImpl(Db *gorm.DB) UsersRepository {
	return &UsersRepositoryImpl{Db: Db}
}

func (u *UsersRepositoryImpl) Delete(usersId int) {
	var users model.Users
	result := u.Db.Where("id = ?", usersId).Delete(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindAll() []model.Users {
	var users []model.Users
	results := u.Db.Find(&users)
	helper.ErrorPanic(results.Error)
	return users
}

func (u *UsersRepositoryImpl) FindById(usersId int) (model.Users, error) {
	var users model.Users
	result := u.Db.Find(&users, usersId)
	if result != nil {
		return users, nil
	} else {
		return users, errors.New("users is not found")
	}
}

func (u *UsersRepositoryImpl) Save(users model.Users) {
	result := u.Db.Create(&users)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) Update(users model.Users) {
	var updateUsers = request.UpdateUsersRequest{
		Id:       users.Id,
		Username: users.Username,
		Email:    users.Email,
		Password: users.Password,
	}
	result := u.Db.Model(&users).Updates(updateUsers)
	helper.ErrorPanic(result.Error)
}

func (u *UsersRepositoryImpl) FindByUsername(username string) (model.Users, error) {
	var users model.Users
	result := u.Db.First(&users, "username = ?", username)

	if result.Error != nil {
		return users, errors.New("invalid username or Password")
	}
	return users, nil
}

func (u *UsersRepositoryImpl) UpdateUser(userID int, updatedUser *model.Users) error {
	var user model.Users
	result := u.Db.First(&user, userID)
	if result.Error != nil {
		return result.Error
	}

	user.Username = updatedUser.Username
	user.Email = updatedUser.Email
	user.Password = updatedUser.Password
	user.IsAdmin = updatedUser.IsAdmin

	result = u.Db.Save(&user)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (u *UsersRepositoryImpl) DeleteUser(userID int) error {
	result := u.Db.Delete(&model.Users{}, userID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
func (repo *UsersRepositoryImpl) GetUserByID(userID int) (*model.Users, error) {
	var user model.Users
	result := repo.Db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}
