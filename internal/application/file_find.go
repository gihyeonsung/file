package application

import "github.com/gihyeonsung/file/internal/domain"

type FileFind struct {
	fileRepository domain.FileRepository
}

func NewFileFind(fileRepository domain.FileRepository) *FileFind {
	return &FileFind{fileRepository: fileRepository}
}

func (u *FileFind) Execute(criteria domain.FileRepositoryCriteria) (domain.FileRepositoryResult, error) {
	return u.fileRepository.Find(criteria)
}
