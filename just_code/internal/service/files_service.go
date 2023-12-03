package service

import (
	"golang-jwttoken/internal/model"
)

type FileService interface {
	SaveFile(file *model.FileModel) error
	GetFileByID(id int) (*model.FileModel, error)
	GetUserByID(userID int) (*model.Users, error)
	GetFilesByUserID(userID int, isAdmin string) ([]model.FileModel, error)
}
