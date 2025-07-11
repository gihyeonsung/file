package application

import (
	"errors"

	"github.com/gihyeonsung/file/internal/domain"
)

type FileDownload struct {
	fileRepository domain.FileRepository
	fileService    FileService
}

func NewFileDownload(fileRepository domain.FileRepository, fileService FileService) *FileDownload {
	return &FileDownload{fileRepository: fileRepository, fileService: fileService}
}

func (u *FileDownload) Execute(id string) ([]byte, error) {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return nil, err
	}

	if file == nil {
		return nil, errors.New("file not found")
	}

	if file.PathRemote == nil {
		return nil, errors.New("file not found")
	}

	bytes, err := u.fileService.Read(*file.PathRemote)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}
