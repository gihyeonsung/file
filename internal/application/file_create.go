package application

import (
	"github.com/gihyeonsung/file/internal/domain"
)

type FileCreate struct {
	fileRepository domain.FileRepository
}

func NewFileCreate(fileRepository domain.FileRepository) *FileCreate {
	return &FileCreate{fileRepository: fileRepository}
}

func (u *FileCreate) Execute(path string) error {
	file, err := domain.NewFile(path)
	if err != nil {
		return err
	}

	return u.fileRepository.Save(file)
}
