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
	result := repo.DB.Preload("User").First(&file, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &file, nil
}

func (repo *FileRepositoryImpl) FindAll(files *[]model.FileModel) error {
	result := repo.DB.Preload("User").Find(files)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *FileRepositoryImpl) FindByUserID(userID int, files *[]model.FileModel) error {
	result := repo.DB.Preload("User").Where("user_id = ?", userID).Find(files)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
