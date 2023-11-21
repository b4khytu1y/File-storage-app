package service

import (
	"awesome/image-storage-service/service/image-storage/internal/entity"

	"gorm.io/gorm"
)

type ImageService struct {
	DB *gorm.DB
}

func NewImageService(db *gorm.DB) *ImageService {
	return &ImageService{
		DB: db,
	}
}

func (service *ImageService) SaveImage(photo *entity.Photo) error {
	result := service.DB.Create(photo)
	return result.Error
}

func (service *ImageService) GetImageByName(name string) (*entity.Photo, error) {
	var photo entity.Photo
	result := service.DB.Where("name = ?", name).First(&photo)
	return &photo, result.Error
}

func (service *ImageService) GetImagesByUserID(userID uint) ([]entity.Photo, error) {
	var photos []entity.Photo
	result := service.DB.Where("user_id = ?", userID).Find(&photos)
	return photos, result.Error
}
