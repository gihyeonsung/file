package application

import "github.com/gihyeonsung/file/internal/domain"

type FileUpload struct {
	fileRepository domain.FileRepository
	fileService    FileService
}

func NewFileUpload(fileRepository domain.FileRepository, fileService FileService) *FileUpload {
	return &FileUpload{fileRepository: fileRepository, fileService: fileService}
}

func (u *FileUpload) Execute(id string, data []byte) error {
	file, err := u.fileRepository.FindOne(id)
	if err != nil {
		return err
	}

	err = u.fileService.Write(file.Path, data)
	if err != nil {
		return err
	}

	pathRemote := file.Path
	size := len(data)
	file.Upload(pathRemote, size)

	return u.fileRepository.Save(file)
}
