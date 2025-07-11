package application

import (
	"errors"

	"github.com/gihyeonsung/file/internal/domain"
)

type FileDelete struct {
	fileRepository domain.FileRepository
	fileService    FileService
}

func NewFileDelete(fileRepository domain.FileRepository, fileService FileService) *FileDelete {
	return &FileDelete{fileRepository: fileRepository, fileService: fileService}
}

func (u *FileDelete) Execute(id string) error {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return err
	}

	if file == nil {
		return errors.New("file not found")
	}

	if file.PathRemote != nil {
		err = u.fileService.Delete(*file.PathRemote)
		if err != nil {
			return err
		}
	}

	return u.fileRepository.Delete(id)
}
