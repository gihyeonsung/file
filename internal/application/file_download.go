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

func (u *FileDownload) Execute(id string) (io.ReadCloser, *string, *int64, error) {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return nil, nil, nil, err
	}

	if file == nil {
		return nil, nil, nil, errors.New("no file")
	}

	if file.PathRemote == nil {
		return nil, nil, nil, errors.New("no file remote path")
	}

	if file.Size == nil {
		return nil, nil, nil, errors.New("no file size")
	}

	r, err := u.fileService.Read(*file.PathRemote)
	if err != nil {
		return nil, nil, nil, err
	}

	return r, file.MimeType, file.Size, nil
}
