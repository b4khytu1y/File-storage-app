package service

import (
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
)

type FileServiceImpl struct {
	fileRepo repository.FileRepository
	userRepo repository.UsersRepository
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
func (s *FileServiceImpl) GetUserByID(userID int) (*model.Users, error) {
	return s.userRepo.GetUserByID(userID)
}
func (s *FileServiceImpl) GetFilesByUserID(userID int) ([]model.FileModel, error) {
	var files []model.FileModel
	err := s.fileRepo.FindByUserID(userID, &files)
	if err != nil {
		return nil, err
	}
	return files, nil
}
