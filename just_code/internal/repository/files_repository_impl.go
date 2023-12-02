package repository

import (
	"golang-jwttoken/internal/model"

	"gorm.io/gorm"
)

type FileRepositoryImpl struct {
	DB *gorm.DB
}

func NewFileRepositoryImpl(db *gorm.DB) FileRepository {
	return &FileRepositoryImpl{DB: db}
}

func (repo *FileRepositoryImpl) Save(file *model.FileModel) error {
	return repo.DB.Create(file).Error
}

func (repo *FileRepositoryImpl) FindByID(id int) (*model.FileModel, error) {
	var file model.FileModel
	err := repo.DB.Where("id = ?", id).First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}
