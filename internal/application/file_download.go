package application

import (
	"errors"
	"io"

	"github.com/gihyeonsung/file/internal/domain"
)

type FileDownload struct {
	fileRepository domain.FileRepository
	fileService    FileService
}

func NewFileDownload(fileRepository domain.FileRepository, fileService FileService) *FileDownload {
	return &FileDownload{fileRepository: fileRepository, fileService: fileService}
}

func (u *FileDownload) Execute(id string) (io.ReadCloser, *string, error) {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return nil, nil, err
	}

	if file == nil {
		return nil, nil, errors.New("no metadata")
	}

	if file.PathRemote == nil {
		return nil, nil, errors.New("no remote path")
	}

	r, err := u.fileService.Read(*file.PathRemote)
	if err != nil {
		return nil, nil, err
	}

	return r, file.MimeType, nil
}
