package repository

import (
	"golang-jwttoken/internal/model"
)

type FileRepository interface {
	Save(file *model.FileModel) error
	FindByID(id int) (*model.FileModel, error)
	FindAll(files *[]model.FileModel) error
	FindByUserID(userID int, files *[]model.FileModel) error
}
