package repository

import (
	"golang-jwttoken/internal/model"
)

type FileRepository interface {
	Save(file *model.FileModel) error
	FindByID(id int) (*model.FileModel, error)
}
