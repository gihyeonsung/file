package application

import (
	"errors"
	"path/filepath"

	"github.com/gihyeonsung/file/internal/domain"
)

type FileCreate struct {
	fileRepository domain.FileRepository
}

func NewFileCreate(fileRepository domain.FileRepository) *FileCreate {
	return &FileCreate{fileRepository: fileRepository}
}

func (u *FileCreate) Execute(path string) error {
	pathNormalized := filepath.Clean(path)
	if !filepath.IsAbs(pathNormalized) {
		return errors.New("not absolute: " + pathNormalized)
	}

	pathName := filepath.Base(pathNormalized)
	if pathName == "" || pathName == "." || pathName == ".." {
		return errors.New("path name: " + pathName)
	}

	result, err := u.fileRepository.Find(&domain.FileRepositoryCriteria{
		Paths: []string{pathNormalized},
	})
	if err != nil {
		return err
	}

	if len(result.Files) > 0 {
		return errors.New("path already exists")
	}

	file, err := domain.NewFile(pathNormalized)
	if err != nil {
		return err
	}

	return u.fileRepository.Save(file)
}
