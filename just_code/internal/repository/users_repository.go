package repository

import "golang-jwttoken/internal/model"

type UsersRepository interface {
	Save(users model.Users)
	Update(users model.Users)
	Delete(usersId int)
	FindById(usersId int) (model.Users, error)
	FindAll() []model.Users
	FindByUsername(username string) (model.Users, error)
	UpdateUser(userID int, updatedUser *model.Users) error
	DeleteUser(userID int) error
	GetUserByID(userID int) (*model.Users, error)
}
