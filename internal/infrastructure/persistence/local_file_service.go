package persistence

import (
	"os"
	"path/filepath"
)

type LocalFileService struct {
	pathBase string
}

func NewLocalFileService(pathBase string) *LocalFileService {
	return &LocalFileService{pathBase: pathBase}
}

func (s *LocalFileService) Write(path string, data []byte) error {
	path = filepath.Join(s.pathBase, path)

	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *LocalFileService) Read(path string) ([]byte, error) {
	path = filepath.Join(s.pathBase, path)

	return os.ReadFile(path)
}

func (s *LocalFileService) Delete(path string) error {
	path = filepath.Join(s.pathBase, path)

	return os.Remove(path)
}
