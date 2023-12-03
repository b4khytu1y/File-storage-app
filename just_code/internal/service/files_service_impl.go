package service

import (
	"golang-jwttoken/internal/model"
	"golang-jwttoken/internal/repository"
)

type FileServiceImpl struct {
	fileRepo repository.FileRepository
	userRepo repository.UsersRepository
}

func NewFileServiceImpl(fileRepo repository.FileRepository, userRepo repository.UsersRepository) FileService {
	return &FileServiceImpl{
		fileRepo: fileRepo,
		userRepo: userRepo,
	}
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
func (s *FileServiceImpl) GetFilesByUserID(userID int, isAdmin string) ([]model.FileModel, error) {
	if isAdmin == "1" {
		var allFiles []model.FileModel
		err := s.fileRepo.FindAll(&allFiles)
		if err != nil {
			return nil, err
		}
		return allFiles, nil
	} else {
		var userFiles []model.FileModel
		err := s.fileRepo.FindByUserID(userID, &userFiles)
		if err != nil {
			return nil, err
		}
		return userFiles, nil
	}
}
