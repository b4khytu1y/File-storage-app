package service

import (
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
)

type FileServiceImpl struct {
	fileRepo repository.FileRepository
}

func NewFileServiceImpl(fileRepo repository.FileRepository) FileService {
	return &FileServiceImpl{fileRepo: fileRepo}
}

func (s *FileServiceImpl) SaveFile(file *model.FileModel) error {
	return s.fileRepo.Save(file)
}

func (s *FileServiceImpl) GetFileByID(id int) (*model.FileModel, error) {
	return s.fileRepo.FindByID(id)
}
