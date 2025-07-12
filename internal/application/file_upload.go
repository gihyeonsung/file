package application

import (
	"io"

	"github.com/gihyeonsung/file/internal/domain"
)

type FileUpload struct {
	fileRepository domain.FileRepository
	fileService    FileService
}

func NewFileUpload(fileRepository domain.FileRepository, fileService FileService) *FileUpload {
	return &FileUpload{fileRepository: fileRepository, fileService: fileService}
}

func (u *FileUpload) Execute(id string, r io.Reader, mimeType string) error {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return err
	}

	size, err := u.fileService.Write(file.Path, r)
	if err != nil {
		return err
	}

	pathRemote := file.Path
	file.Upload(pathRemote, size, mimeType)

	return u.fileRepository.Save(file)
}
