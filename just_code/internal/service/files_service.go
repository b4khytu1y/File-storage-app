package service

import (
	"golang-jwttoken/internal/model"
)

type FileService interface {
	SaveFile(file *model.FileModel) error
	GetFileByID(id int) (*model.FileModel, error)
	GetUserByID(id int) (*model.Users, error)
}
